package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	client *ent.Client
}

func New(conn string) (*Repo, error) {
	db, err := sql.Open("pgx", conn)
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

func (r *Repo) GetRecords(ctx context.Context,
	metadataPrefix string,
	set string,
	from string,
	until string,
) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	return r.getRecords(ctx, metadataPrefix, set, from, until, 0)
}

func (r *Repo) GetMoreRecords(ctx context.Context, tokenValue string) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	// TODO validate token
	// TODO encrypt token
	tokenParts := strings.SplitN(tokenValue, "|", 5)
	metadataPrefix := tokenParts[0]
	setSpec := tokenParts[1]
	from := tokenParts[2]
	until := tokenParts[3]
	lastID, err := strconv.ParseInt(tokenParts[4], 10, 64)
	if err != nil {
		return nil, nil, err
	}

	return r.getRecords(ctx, metadataPrefix, setSpec, from, until, lastID)
}

// TODO set
func (r *Repo) getRecords(ctx context.Context, metadataPrefix, setSpec, from, until string, lastID int64) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	where := []predicate.Record{
		record.HasMetadataFormatWith(metadataformat.PrefixEQ(metadataPrefix)),
	}
	if setSpec != "" {
		where = append(where, record.HasSetsWith(set.SpecEQ(setSpec)))
	}
	if from != "" {
		dt, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return nil, nil, err
		}
		where = append(where, record.DatestampGTE(dt))
	}
	if until != "" {
		dt, err := time.Parse(time.RFC3339, until)
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

	if lastID > 0 {
		where = append(where, record.IDGT(lastID))
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
		token = &oaipmh.ResumptionToken{
			CompleteListSize: n,
			Value:            fmt.Sprintf("%s|%s|%s|%s|%d", metadataPrefix, setSpec, from, until, rows[len(rows)-1].ID),
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

// TODO sets
func (r *Repo) AddRecord(ctx context.Context, identifier, metadataPrefix, metadata string, setSpecs []string) error {
	formatID, err := r.client.MetadataFormat.Query().
		Where(metadataformat.PrefixEQ(metadataPrefix)).OnlyID(ctx)
	if err != nil {
		return err
	}

	return r.client.Record.Create().
		SetIdentifier(identifier).
		SetMetadataFormatID(formatID).
		SetMetadata(metadata).
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
