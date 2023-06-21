package repository

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	sqldialect "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/oai-service/ent"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/migrate"
	"github.com/ugent-library/oai-service/ent/record"
	"github.com/ugent-library/oai-service/ent/set"
	"github.com/ugent-library/oai-service/models"
)

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

func (r *Repo) GetAllMetadataFormats(ctx context.Context) ([]*models.MetadataFormat, error) {
	rows, err := r.client.MetadataFormat.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	formats := make([]*models.MetadataFormat, len(rows))
	for i, row := range rows {
		formats[i] = &models.MetadataFormat{
			Prefix:    row.Prefix,
			Schema:    row.Schema,
			Namespace: row.Namespace,
		}
	}
	return formats, nil
}

func (r *Repo) HasRecord(ctx context.Context, identifier string) (bool, error) {
	return r.client.Record.Query().
		Where(record.IdentifierEQ(identifier)).
		Exist(ctx)
}

func (r *Repo) GetRecord(ctx context.Context, identifier, metadataPrefix string) (*models.Record, error) {
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
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rec := &models.Record{
		Identifier: row.Identifier,
		Datestamp:  row.Datestamp,
		Deleted:    row.Deleted,
	}
	for _, set := range row.Edges.Sets {
		rec.SetSpecs = append(rec.SetSpecs, set.Spec)
	}
	if !row.Deleted {
		rec.Metadata = row.Metadata
	}

	return rec, nil
}
