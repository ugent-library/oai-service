package repositories

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	sqldialect "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/crypt"

	"github.com/ugent-library/oai-service/ent"
	"github.com/ugent-library/oai-service/ent/item"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/migrate"
	"github.com/ugent-library/oai-service/ent/predicate"
	"github.com/ugent-library/oai-service/ent/record"
	"github.com/ugent-library/oai-service/ent/set"
	"github.com/ugent-library/oai-service/oaipmh"
)

var ErrNotFound = errors.New("not found")

type Repo struct {
	client *ent.Client
	config Config
}

type Config struct {
	Conn   string
	Secret []byte
}

type setCursor struct {
	LastID string `json:"l"`
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
		Where(metadataformat.IDEQ(prefix)).
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
			MetadataPrefix:    row.ID,
			Schema:            row.Schema,
			MetadataNamespace: row.Namespace,
		}
	}
	return formats, nil
}

func (r *Repo) AddMetadataFormat(ctx context.Context, prefix, schema, namespace string) error {
	return r.client.MetadataFormat.Create().
		SetID(prefix).
		SetSchema(schema).
		SetNamespace(namespace).
		OnConflictColumns(metadataformat.FieldID).
		UpdateNewValues().
		Exec(ctx)
}

func (r *Repo) HasSets(ctx context.Context) (bool, error) {
	return r.client.Set.Query().Exist(ctx)
}

func (r *Repo) HasSet(ctx context.Context, spec string) (bool, error) {
	return r.client.Set.Query().
		Where(set.IDEQ(spec)).
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
	total, err := r.client.Set.Query().
		Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	if total == 0 {
		return nil, nil, nil
	}

	rows, err := r.client.Set.Query().
		Where(set.IDGT(c.LastID)).
		Order(ent.Asc(set.FieldID)).
		Limit(100).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	sets := make([]*oaipmh.Set, len(rows))
	for i, row := range rows {
		sets[i] = &oaipmh.Set{
			SetSpec: row.ID,
			SetName: row.Name,
			SetDescription: &oaipmh.Payload{
				Content: row.Description,
			},
		}
	}

	var token *oaipmh.ResumptionToken
	if total > len(rows) {
		tokenValue, err := r.encodeCursor(setCursor{
			LastID: rows[len(rows)-1].ID,
		})
		if err != nil {
			return nil, nil, err
		}
		token = &oaipmh.ResumptionToken{
			CompleteListSize: total,
			Value:            tokenValue,
		}
	}

	return sets, token, nil
}

func (r *Repo) AddSet(ctx context.Context, spec, name, description string) error {
	return r.client.Set.Create().
		SetID(spec).
		SetName(name).
		SetDescription(description).
		OnConflictColumns(set.FieldID).
		UpdateNewValues().
		Exec(ctx)
}

func (r *Repo) GetEarliestDatestamp(ctx context.Context) (time.Time, error) {
	row, err := r.client.Record.Query().
		Select(record.FieldDatestamp).
		Order(ent.Asc(record.FieldDatestamp)).
		First(ctx)
	if ent.IsNotFound(err) {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, err
	}
	return row.Datestamp, nil
}

func (r *Repo) HasRecord(ctx context.Context, identifier string) (bool, error) {
	return r.client.Record.Query().
		Where(record.ItemIDEQ(identifier)).
		Exist(ctx)
}

func (r *Repo) GetRecord(ctx context.Context, identifier, prefix string) (*oaipmh.Record, error) {
	row, err := r.client.Record.Query().
		Where(
			record.ItemIDEQ(identifier),
			record.MetadataFormatIDEQ(prefix),
		).
		WithItem(func(q *ent.ItemQuery) {
			q.WithSets(func(q *ent.SetQuery) {
				q.Select(set.FieldID)
			})
		}).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, oaipmh.ErrCannotDisseminateFormat
	}
	if err != nil {
		return nil, err
	}

	rec := &oaipmh.Record{
		Header: &oaipmh.Header{
			Identifier: row.ItemID,
			Datestamp:  row.Datestamp.UTC().Format(time.RFC3339),
		},
	}
	for _, set := range row.Edges.Item.Edges.Sets {
		rec.Header.SetSpecs = append(rec.Header.SetSpecs, set.ID)
	}
	if row.Metadata == nil {
		rec.Header.Status = "deleted"
	} else {
		rec.Metadata = &oaipmh.Payload{
			Content: *row.Metadata,
		}
	}

	return rec, nil
}

// TODO this loads the complete record, make an efficient version
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
		record.MetadataFormatIDEQ(c.MetadataPrefix),
	}
	if c.SetSpec != "" {
		where = append(where, record.HasItemWith(
			item.HasSetsWith(set.IDEQ(c.SetSpec)),
		))
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

	total, err := r.client.Record.Query().
		Where(where...).
		Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	if total == 0 {
		return nil, nil, nil
	}

	if c.LastID > 0 {
		where = append(where, record.IDGT(c.LastID))
	}

	rows, err := r.client.Record.Query().
		Where(where...).
		WithItem(func(q *ent.ItemQuery) {
			q.WithSets(func(q *ent.SetQuery) {
				q.Select(set.FieldID)
			})
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
				Identifier: row.ItemID,
				Datestamp:  row.Datestamp.UTC().Format(time.RFC3339),
			},
		}
		for _, s := range row.Edges.Item.Edges.Sets {
			rec.Header.SetSpecs = append(rec.Header.SetSpecs, s.ID)
		}
		if row.Metadata == nil {
			rec.Header.Status = "deleted"
		} else {
			rec.Metadata = &oaipmh.Payload{
				Content: *row.Metadata,
			}
		}
		recs[i] = rec
	}

	var token *oaipmh.ResumptionToken
	if total > len(rows) {
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
			CompleteListSize: total,
			Value:            tokenValue,
		}
	}

	return recs, token, nil
}

// TODO scan directly into []*oaipmh.MetadataFormat?
func (r *Repo) GetRecordMetadataFormats(ctx context.Context, identifier string) ([]*oaipmh.MetadataFormat, error) {
	rows, err := r.client.Record.Query().
		Where(record.ItemIDEQ(identifier)).
		QueryMetadataFormat().All(ctx)
	if err != nil {
		return nil, err
	}
	formats := make([]*oaipmh.MetadataFormat, len(rows))
	for i, row := range rows {
		formats[i] = &oaipmh.MetadataFormat{
			MetadataPrefix:    row.ID,
			Schema:            row.Schema,
			MetadataNamespace: row.Namespace,
		}
	}
	return formats, nil
}

func (r *Repo) AddItem(ctx context.Context, identifier string, specs []string) error {
	sql := `
  	with add_item as (
			insert into items (id) values ($1)
	    on conflict (id)
		  do nothing
	  ), del_sets as (
	  	delete from item_sets
	  	where item_id = $1 and set_id != all($2)
	  )
    insert into item_sets (item_id, set_id)
	  values ($1, unnest($2))
	  on conflict (item_id, set_id)
	  do nothing
	`
	_, err := r.client.ExecContext(ctx, sql, identifier, specs)
	return err
}

func (r *Repo) AddRecord(ctx context.Context, identifier, prefix, content string) error {
	sql := `
  	with add_item as (
			insert into items (id) values ($1)
	    on conflict (id)
		  do nothing
	  )
    insert into records (item_id, metadata_format_id, metadata, datestamp)
	  values ($1, $2, $3, current_timestamp)
	  on conflict (item_id, metadata_format_id)
	  do update set metadata = excluded.metadata, datestamp = excluded.datestamp
	  where records.metadata != excluded.metadata
	`
	_, err := r.client.ExecContext(ctx, sql, identifier, prefix, content)
	return err
}

// TODO add prefix argument
func (r *Repo) DeleteRecord(ctx context.Context, identifier, prefix string) error {
	return r.client.Record.Update().
		Where(
			record.MetadataFormatIDEQ(prefix),
			record.ItemIDEQ(identifier),
			record.MetadataNotNil(),
		).
		ClearMetadata().
		SetDatestamp(time.Now()).
		Exec(ctx)
}

func (r *Repo) encodeCursor(c any) (string, error) {
	plaintext, _ := json.Marshal(c)
	ciphertext, err := crypt.Encrypt(r.config.Secret, plaintext)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (r *Repo) decodeCursor(encryptedCursor string, c any) error {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedCursor)
	if err != nil {
		return err
	}
	plaintext, err := crypt.Decrypt(r.config.Secret, ciphertext)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, c)
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
