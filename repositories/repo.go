package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"entgo.io/ent/dialect"
	sqldialect "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ugent-library/crypt"

	"github.com/ugent-library/oai-service/ent"
	"github.com/ugent-library/oai-service/ent/item"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/predicate"
	"github.com/ugent-library/oai-service/ent/record"
	"github.com/ugent-library/oai-service/ent/set"
	"github.com/ugent-library/oai-service/oaipmh"
)

var ErrNotFound = errors.New("not found")

const PageLimit = 100

type Repo struct {
	client *ent.Client
	config Config
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

	// err = client.Schema.Create(context.TODO(),
	// 	migrate.WithDropIndex(true),
	// )
	// if err != nil {
	// 	return nil, err
	// }

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

func (r *Repo) AddMetadataFormat(ctx context.Context, prefix, schema, namespace string) error {
	return r.client.MetadataFormat.Create().
		SetPrefix(prefix).
		SetSchema(schema).
		SetNamespace(namespace).
		OnConflictColumns(metadataformat.FieldPrefix).
		UpdateNewValues().
		Exec(ctx)
}

func (r *Repo) HasSets(ctx context.Context) (bool, error) {
	return r.client.Set.Query().Exist(ctx)
}

func (r *Repo) HasSet(ctx context.Context, spec string) (bool, error) {
	return r.client.Set.Query().
		Where(set.Or(set.SpecEQ(spec), set.SpecHasPrefix(spec+":"))).
		Exist(ctx)
}

func (r *Repo) GetSets(ctx context.Context) ([]*oaipmh.Set, *oaipmh.ResumptionToken, error) {
	return r.getSets(ctx, setCursor{})
}

func (r *Repo) GetMoreSets(ctx context.Context, tokenValue string) ([]*oaipmh.Set, *oaipmh.ResumptionToken, error) {
	c := setCursor{}
	if err := crypt.DecryptValue(r.config.Secret, tokenValue, &c); err != nil {
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
		Limit(PageLimit).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	sets := make([]*oaipmh.Set, len(rows))
	for i, row := range rows {
		sets[i] = &oaipmh.Set{
			SetSpec: row.Spec,
			SetName: row.Name,
			SetDescription: &oaipmh.Payload{
				XML: row.Description,
			},
		}
	}

	var token *oaipmh.ResumptionToken
	if len(rows) >= PageLimit {
		tokenValue, err := crypt.EncryptValue(r.config.Secret, setCursor{
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
		SetSpec(spec).
		SetName(name).
		SetDescription(description).
		OnConflictColumns(set.FieldSpec).
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
		Where(record.HasItemWith(item.IdentifierEQ(identifier))).
		Exist(ctx)
}

func (r *Repo) GetRecord(ctx context.Context, identifier, prefix string) (*oaipmh.Record, error) {
	row, err := r.client.Record.Query().
		Where(
			record.HasItemWith(item.IdentifierEQ(identifier)),
			record.HasMetadataFormatWith(metadataformat.PrefixEQ(prefix)),
		).
		WithItem(func(q *ent.ItemQuery) {
			q.WithSets(func(q *ent.SetQuery) {
				q.Select(set.FieldSpec)
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
			Identifier: row.Edges.Item.Identifier,
			Datestamp:  row.Datestamp.UTC().Format(time.RFC3339),
		},
	}
	for _, set := range row.Edges.Item.Edges.Sets {
		rec.Header.SetSpecs = append(rec.Header.SetSpecs, set.Spec)
	}
	if row.Metadata == nil {
		rec.Header.Status = "deleted"
	} else {
		rec.Metadata = &oaipmh.Payload{
			XML: *row.Metadata,
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
	if err := crypt.DecryptValue(r.config.Secret, tokenValue, &c); err != nil {
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
	if err := crypt.DecryptValue(r.config.Secret, tokenValue, &c); err != nil {
		return nil, nil, err
	}
	return r.getRecords(ctx, c)
}

func (r *Repo) getRecords(ctx context.Context, c recordCursor) ([]*oaipmh.Record, *oaipmh.ResumptionToken, error) {
	where := []predicate.Record{
		record.HasMetadataFormatWith(metadataformat.PrefixEQ(c.MetadataPrefix)),
	}
	if c.SetSpec != "" {
		where = append(where, record.HasItemWith(
			item.HasSetsWith(set.Or(set.SpecEQ(c.SetSpec), set.SpecHasPrefix(c.SetSpec+":"))),
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
				q.Select(set.FieldSpec)
			})
		}).
		Order(ent.Asc(record.FieldID)).
		Limit(PageLimit).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}
	recs := make([]*oaipmh.Record, len(rows))
	for i, row := range rows {
		rec := &oaipmh.Record{
			Header: &oaipmh.Header{
				Identifier: row.Edges.Item.Identifier,
				Datestamp:  row.Datestamp.UTC().Format(time.RFC3339),
			},
		}
		for _, s := range row.Edges.Item.Edges.Sets {
			rec.Header.SetSpecs = append(rec.Header.SetSpecs, s.Spec)
		}
		if row.Metadata == nil {
			rec.Header.Status = "deleted"
		} else {
			rec.Metadata = &oaipmh.Payload{
				XML: *row.Metadata,
			}
		}
		recs[i] = rec
	}

	var token *oaipmh.ResumptionToken
	if len(rows) >= PageLimit {
		tokenValue, err := crypt.EncryptValue(r.config.Secret, recordCursor{
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
		Where(record.HasItemWith(item.IdentifierEQ(identifier))).
		QueryMetadataFormat().All(ctx)
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

func (r *Repo) AddItem(ctx context.Context, identifier string, specs []string) error {
	sql := `
  	with add_item as (
	  insert into items (identifier) values($1)
	  on conflict (identifier)
	    do nothing
	  returning id
	), item as (
      select id from add_item
      union
      select id from items where identifier = $1
	), new_sets as (
	  select id from sets where spec = any($2)
	), del_sets as (
	  delete from item_sets
	  using item, new_sets
	  where item_id = item.id and set_id not in (select id from new_sets)
	)
    insert into item_sets (item_id, set_id)
	  select item.id, new_sets.id
	  from item, new_sets
	  on conflict (item_id, set_id)
	    do nothing
	`
	_, err := r.client.ExecContext(ctx, sql, identifier, specs)
	return err
}

func (r *Repo) AddRecord(ctx context.Context, identifier, prefix, metadata string) error {
	sql := `
  	with add_item as (
	  insert into items (identifier) values($1)
	  on conflict (identifier)
	    do nothing
	  returning id
	), item as (
      select id from add_item
      union
      select id from items where identifier = $1
	), fmt as (
	  select id from metadata_formats where prefix = $2
	)
    insert into records (item_id, metadata_format_id, metadata, datestamp)
	  select item.id, fmt.id, $3, current_timestamp
	  from item, fmt
	  on conflict (item_id, metadata_format_id)
	    do update set metadata = excluded.metadata, datestamp = excluded.datestamp
	    where records.metadata != excluded.metadata
	`
	_, err := r.client.ExecContext(ctx, sql, identifier, prefix, metadata)
	return err
}

func (r *Repo) DeleteRecord(ctx context.Context, identifier, prefix string) error {
	sql := `
  	with add_item as (
	  insert into items (identifier) values($1)
	  on conflict (identifier)
		do nothing
	  returning id
	), item as (
      select id from add_item
      union
      select id from items where identifier = $1
	), fmt as (
	  select id from metadata_formats where prefix = $2
	)
    insert into records (item_id, metadata_format_id, metadata, datestamp)
	  select item.id, fmt.id, NULL, current_timestamp
	  from item, fmt
	  on conflict (item_id, metadata_format_id)
	  	do update set metadata = NULL, datestamp = excluded.datestamp
	  	where records.metadata is not null
	`
	_, err := r.client.ExecContext(ctx, sql, identifier, prefix)
	return err
}
