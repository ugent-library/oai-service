package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	sqldialect "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/oai-service/ent"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/migrate"
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

func (r *Repo) HasRecord(ctx context.Context, identifier string) (bool, error) {
	return r.client.Record.Query().
		Where(record.IdentifierEQ(identifier)).
		Exist(ctx)
}

func (r *Repo) GetRecords(ctx context.Context, metadataPrefix string) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	// TODO ent can't do count and select in one query
	n, err := r.client.Record.Query().
		Where(
			record.HasMetadataFormatWith(metadataformat.PrefixEQ(metadataPrefix)),
		).
		Count(ctx)
	if err != nil {
		return nil, nil, err
	}
	if n == 0 {
		return nil, nil, nil
	}

	rows, err := r.client.Record.Query().
		Where(
			record.HasMetadataFormatWith(metadataformat.PrefixEQ(metadataPrefix)),
		).
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
			Value:            fmt.Sprint(rows[len(rows)-1].ID),
		}
	}

	return recs, token, nil
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
