package repository

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"

	"entgo.io/ent/dialect"
	sqldialect "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/oai-service/ent"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/migrate"
	"github.com/ugent-library/oai-service/ent/predicate"
	"github.com/ugent-library/oai-service/ent/record"
	"github.com/ugent-library/oai-service/ent/set"
	"github.com/ugent-library/oai-service/oaipmh"
)

var ErrNotFound = errors.New("not found")

type Repo struct {
	config Config
	client *ent.Client
}

type Config struct {
	Conn   string
	Secret []byte
}

type setCursor struct {
	LastID int64 `json:"l"`
}

type recordCursor struct {
	MetadataPrefix string `json:"m"`
	SetSpec        string `json:"s"`
	From           string `json:"f"`
	Until          string `json:"u"`
	LastID         int64  `json:"l"`
}

func New(c Config) (*Repo, error) {
	db, err := sql.Open("pgx", c.Conn)
	if err != nil {
		return nil, err
	}

	driver := sqldialect.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver))

	err = client.Schema.Create(context.TODO(),
		migrate.WithDropIndex(true),
	)
	if err != nil {
		return nil, err
	}

	return &Repo{
		config: c,
		client: client,
	}, nil
}

func (r *Repo) HasMetadataFormat(ctx context.Context, prefix string) (bool, error) {
	return r.client.MetadataFormat.Query().
		Where(metadataformat.PrefixEQ(prefix)).
		Exist(ctx)
}

func (r *Repo) GetMetadataFormats(ctx context.Context) ([]*oaipmh.MetadataFormat, error) {
	rows, err := r.client.MetadataFormat.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	formats := make([]*oaipmh.MetadataFormat, len(rows))
	for i, row := range rows {
		formats[i] = &oaipmh.MetadataFormat{
			MetadataPrefix:    row.Prefix,
			Schema:            row.Schema,
			MetadataNamespace: row.Namespace,
		}
	}
	return formats, nil
}

func (r *Repo) HasSet(ctx context.Context, spec string) (bool, error) {
	return r.client.Set.Query().
		Where(set.SpecEQ(spec)).
		Exist(ctx)
}

func (r *Repo) GetSets(ctx context.Context) ([]*oaipmh.Set, *oaipmh.ResumptionToken, error) {
	return r.getSets(ctx, setCursor{})
}

func (r *Repo) GetMoreSets(ctx context.Context, tokenValue string) ([]*oaipmh.Set, *oaipmh.ResumptionToken, error) {
	c := setCursor{}
	if err := r.decodeCursor(tokenValue, &c); err != nil {
		return nil, nil, err
	}
	return r.getSets(ctx, c)
}

func (r *Repo) getSets(ctx context.Context, c setCursor) ([]*oaipmh.Set, *oaipmh.ResumptionToken, error) {
	// TODO ent can't do count and select in one query
	n, err := r.client.Set.Query().
		Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	if n == 0 {
		return nil, nil, nil
	}

	var where []predicate.Set
	if c.LastID > 0 {
		where = append(where, set.IDGT(c.LastID))
	}

	rows, err := r.client.Set.Query().
		Where(where...).
		Order(ent.Asc(set.FieldID)).
		Limit(100).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	sets := make([]*oaipmh.Set, len(rows))
	for i, row := range rows {
		sets[i] = &oaipmh.Set{
			Spec: row.Spec,
			Name: row.Name,
			Description: &oaipmh.Payload{
				XML: row.Description,
			},
		}
	}

	var token *oaipmh.ResumptionToken
	if n > len(rows) {
		tokenValue, err := r.encodeCursor(setCursor{
			LastID: rows[len(rows)-1].ID,
		})
		if err != nil {
			return nil, nil, err
		}
		token = &oaipmh.ResumptionToken{
			CompleteListSize: n,
			Value:            tokenValue,
		}
	}

	return sets, token, nil
}

func (r *Repo) AddSet(ctx context.Context, spec, name, description string) error {
	return r.client.Set.Create().
		SetSpec(spec).
		SetName(name).
		SetDescription(description).
		OnConflictColumns(set.FieldSpec).
		UpdateNewValues().
		Exec(ctx)
}

func (r *Repo) HasRecord(ctx context.Context, identifier string) (bool, error) {
	return r.client.Record.Query().
		Where(record.IdentifierEQ(identifier)).
		Exist(ctx)
}

func (r *Repo) GetRecord(ctx context.Context, identifier, metadataPrefix string) (*oaipmh.Record, error) {
	row, err := r.client.Record.Query().
		Where(
			record.IdentifierEQ(identifier),
			record.HasMetadataFormatWith(metadataformat.PrefixEQ(metadataPrefix)),
		).
		WithSets(func(q *ent.SetQuery) {
			q.Select(set.FieldSpec)
		}).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rec := &oaipmh.Record{
		Header: &oaipmh.Header{
			Identifier: row.Identifier,
			Datestamp:  row.Datestamp.UTC().Format(time.RFC3339),
		},
	}
	for _, set := range row.Edges.Sets {
		rec.Header.SetSpecs = append(rec.Header.SetSpecs, set.Spec)
	}
	if row.Deleted {
		rec.Header.Status = "deleted"
	} else {
		rec.Metadata = &oaipmh.Payload{
			XML: row.Metadata,
		}
	}

	return rec, nil
}

// TODO this loads the complete record, maken an efficient version
func (r *Repo) GetIdentifiers(ctx context.Context,
	metadataPrefix string,
	setSpec string,
	from string,
	until string,
) ([]*oaipmh.Header, *oaipmh.ResumptionToken, error) {
	recs, token, err := r.getRecords(ctx, recordCursor{
		MetadataPrefix: metadataPrefix,
		SetSpec:        setSpec,
		From:           from,
		Until:          until,
	})
	if err != nil {
		return nil, nil, err
	}
	headers := make([]*oaipmh.Header, len(recs))
	for i, rec := range recs {
		headers[i] = rec.Header
	}
	return headers, token, nil
}

// TODO this loads the complete record, maken an efficient version
func (r *Repo) GetMoreIdentifiers(ctx context.Context, tokenValue string) ([]*oaipmh.Header, *oaipmh.ResumptionToken, error) {
	c := recordCursor{}
	if err := r.decodeCursor(tokenValue, &c); err != nil {
		return nil, nil, err
	}
	recs, token, err := r.getRecords(ctx, c)
	if err != nil {
		return nil, nil, err
	}
	headers := make([]*oaipmh.Header, len(recs))
	for i, rec := range recs {
		headers[i] = rec.Header
	}
	return headers, token, nil
}

func (r *Repo) GetRecords(ctx context.Context,
	metadataPrefix string,
	setSpec string,
	from string,
	until string,
) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	return r.getRecords(ctx, recordCursor{
		MetadataPrefix: metadataPrefix,
		SetSpec:        setSpec,
		From:           from,
		Until:          until,
	})
}

func (r *Repo) GetMoreRecords(ctx context.Context, tokenValue string) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	c := recordCursor{}
	if err := r.decodeCursor(tokenValue, &c); err != nil {
		return nil, nil, err
	}
	return r.getRecords(ctx, c)
}

func (r *Repo) getRecords(ctx context.Context, c recordCursor) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	where := []predicate.Record{
		record.HasMetadataFormatWith(metadataformat.PrefixEQ(c.MetadataPrefix)),
	}
	if c.SetSpec != "" {
		where = append(where, record.HasSetsWith(set.SpecEQ(c.SetSpec)))
	}
	if c.From != "" {
		dt, err := time.Parse(time.RFC3339, c.From)
		if err != nil {
			return nil, nil, err
		}
		where = append(where, record.DatestampGTE(dt))
	}
	if c.Until != "" {
		dt, err := time.Parse(time.RFC3339, c.Until)
		if err != nil {
			return nil, nil, err
		}
		where = append(where, record.DatestampLTE(dt))
	}

	// TODO ent can't do count and select in one query
	n, err := r.client.Record.Query().
		Where(where...).
		Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	if n == 0 {
		return nil, nil, nil
	}

	if c.LastID > 0 {
		where = append(where, record.IDGT(c.LastID))
	}

	rows, err := r.client.Record.Query().
		Where(where...).
		WithSets(func(q *ent.SetQuery) {
			q.Select(set.FieldSpec)
		}).
		Order(ent.Asc(record.FieldID)).
		Limit(100).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	recs := make([]*oaipmh.Record, len(rows))
	for i, row := range rows {
		rec := &oaipmh.Record{
			Header: &oaipmh.Header{
				Identifier: row.Identifier,
				Datestamp:  row.Datestamp.UTC().Format(time.RFC3339),
			},
		}
		for _, set := range row.Edges.Sets {
			rec.Header.SetSpecs = append(rec.Header.SetSpecs, set.Spec)
		}
		if row.Deleted {
			rec.Header.Status = "deleted"
		} else {
			rec.Metadata = &oaipmh.Payload{
				XML: row.Metadata,
			}
		}
		recs[i] = rec
	}

	var token *oaipmh.ResumptionToken
	if n > len(rows) {
		tokenValue, err := r.encodeCursor(recordCursor{
			MetadataPrefix: c.MetadataPrefix,
			SetSpec:        c.SetSpec,
			From:           c.From,
			Until:          c.Until,
			LastID:         rows[len(rows)-1].ID,
		})
		if err != nil {
			return nil, nil, err
		}
		token = &oaipmh.ResumptionToken{
			CompleteListSize: n,
			Value:            tokenValue,
		}
	}

	return recs, token, nil
}

func (r *Repo) GetRecordMetadataFormats(ctx context.Context, identifier string) ([]*oaipmh.MetadataFormat, error) {
	rows, err := r.client.Record.Query().
		Where(record.IdentifierEQ(identifier)).
		QueryMetadataFormat().All(ctx)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	formats := make([]*oaipmh.MetadataFormat, len(rows))
	for i, row := range rows {
		formats[i] = &oaipmh.MetadataFormat{
			MetadataPrefix:    row.Prefix,
			Schema:            row.Schema,
			MetadataNamespace: row.Namespace,
		}
	}
	return formats, nil
}

// TODO remove old sets first
func (r *Repo) AddRecord(ctx context.Context, identifier, metadataPrefix, metadata string, setSpecs []string) error {
	formatID, err := r.client.MetadataFormat.Query().
		Where(metadataformat.PrefixEQ(metadataPrefix)).OnlyID(ctx)
	if err != nil {
		return err
	}

	setIDs := make([]int64, len(setSpecs))
	for i, spec := range setSpecs {
		id, err := r.client.Set.Query().
			Where(set.SpecEQ(spec)).
			OnlyID(ctx)
		if err != nil {
			return err
		}
		setIDs[i] = id
	}

	return r.client.Record.Create().
		SetIdentifier(identifier).
		SetMetadataFormatID(formatID).
		SetMetadata(metadata).
		AddSetIDs(setIDs...).
		OnConflictColumns(record.FieldIdentifier, record.FieldMetadataFormatID).
		UpdateNewValues().
		Exec(ctx)
}

func (r *Repo) DeleteRecord(ctx context.Context, identifier string) error {
	return r.client.Record.Update().
		Where(record.IdentifierEQ(identifier)).
		SetDeleted(true).
		SetNillableMetadata(nil).
		Exec(ctx)
}

func (r *Repo) encodeCursor(c any) (string, error) {
	msg, _ := json.Marshal(c)

	// Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(r.config.Secret)
	if err != nil {
		return "", err
	}

	// Wrap the cipher block in Galois Counter Mode.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a unique nonce containing 12 random bytes.
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	// 	// Encrypt the data using aesGCM.Seal(). By passing the nonce as the first
	// 	// parameter, the encrypted message will be appended to the nonce so
	// 	// that the encrypted message will be in the format
	// 	// "{nonce}{encrypted message}".
	cryptedMsg := gcm.Seal(nonce, nonce, msg, nil)

	// Encode as a url safe base64 string.
	return base64.URLEncoding.EncodeToString(cryptedMsg), nil
}

func (r *Repo) decodeCursor(encodedMsg string, c any) error {
	// Decode base64.
	cryptedMsg, err := base64.URLEncoding.DecodeString(encodedMsg)
	if err != nil {
		return err
	}

	// Create a new AES cipher block from the secret key.
	block, err := aes.NewCipher(r.config.Secret)
	if err != nil {
		return err
	}

	// Wrap the cipher block in Galois Counter Mode.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()

	// Avoid potential 'index out of range' panic in the next step.
	if len(cryptedMsg) < nonceSize {
		return oaipmh.ErrBadResumptionToken
	}

	// Split cryptedMsg in nonce and encrypted message and use gcm.Open() to
	// decrypt and authenticate the data.
	msg, err := gcm.Open(nil, cryptedMsg[:nonceSize], cryptedMsg[nonceSize:], nil)
	if err != nil {
		return oaipmh.ErrBadResumptionToken
	}

	err = json.Unmarshal(msg, c)

	return err
}
