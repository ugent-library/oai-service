// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ugent-library/oai-service/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/ugent-library/oai-service/ent/metadata"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/record"
	"github.com/ugent-library/oai-service/ent/set"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Metadata is the client for interacting with the Metadata builders.
	Metadata *MetadataClient
	// MetadataFormat is the client for interacting with the MetadataFormat builders.
	MetadataFormat *MetadataFormatClient
	// Record is the client for interacting with the Record builders.
	Record *RecordClient
	// Set is the client for interacting with the Set builders.
	Set *SetClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Metadata = NewMetadataClient(c.config)
	c.MetadataFormat = NewMetadataFormatClient(c.config)
	c.Record = NewRecordClient(c.config)
	c.Set = NewSetClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:            ctx,
		config:         cfg,
		Metadata:       NewMetadataClient(cfg),
		MetadataFormat: NewMetadataFormatClient(cfg),
		Record:         NewRecordClient(cfg),
		Set:            NewSetClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:            ctx,
		config:         cfg,
		Metadata:       NewMetadataClient(cfg),
		MetadataFormat: NewMetadataFormatClient(cfg),
		Record:         NewRecordClient(cfg),
		Set:            NewSetClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Metadata.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Metadata.Use(hooks...)
	c.MetadataFormat.Use(hooks...)
	c.Record.Use(hooks...)
	c.Set.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Metadata.Intercept(interceptors...)
	c.MetadataFormat.Intercept(interceptors...)
	c.Record.Intercept(interceptors...)
	c.Set.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *MetadataMutation:
		return c.Metadata.mutate(ctx, m)
	case *MetadataFormatMutation:
		return c.MetadataFormat.mutate(ctx, m)
	case *RecordMutation:
		return c.Record.mutate(ctx, m)
	case *SetMutation:
		return c.Set.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// MetadataClient is a client for the Metadata schema.
type MetadataClient struct {
	config
}

// NewMetadataClient returns a client for the Metadata from the given config.
func NewMetadataClient(c config) *MetadataClient {
	return &MetadataClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `metadata.Hooks(f(g(h())))`.
func (c *MetadataClient) Use(hooks ...Hook) {
	c.hooks.Metadata = append(c.hooks.Metadata, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `metadata.Intercept(f(g(h())))`.
func (c *MetadataClient) Intercept(interceptors ...Interceptor) {
	c.inters.Metadata = append(c.inters.Metadata, interceptors...)
}

// Create returns a builder for creating a Metadata entity.
func (c *MetadataClient) Create() *MetadataCreate {
	mutation := newMetadataMutation(c.config, OpCreate)
	return &MetadataCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Metadata entities.
func (c *MetadataClient) CreateBulk(builders ...*MetadataCreate) *MetadataCreateBulk {
	return &MetadataCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Metadata.
func (c *MetadataClient) Update() *MetadataUpdate {
	mutation := newMetadataMutation(c.config, OpUpdate)
	return &MetadataUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *MetadataClient) UpdateOne(m *Metadata) *MetadataUpdateOne {
	mutation := newMetadataMutation(c.config, OpUpdateOne, withMetadata(m))
	return &MetadataUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *MetadataClient) UpdateOneID(id int64) *MetadataUpdateOne {
	mutation := newMetadataMutation(c.config, OpUpdateOne, withMetadataID(id))
	return &MetadataUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Metadata.
func (c *MetadataClient) Delete() *MetadataDelete {
	mutation := newMetadataMutation(c.config, OpDelete)
	return &MetadataDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *MetadataClient) DeleteOne(m *Metadata) *MetadataDeleteOne {
	return c.DeleteOneID(m.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *MetadataClient) DeleteOneID(id int64) *MetadataDeleteOne {
	builder := c.Delete().Where(metadata.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &MetadataDeleteOne{builder}
}

// Query returns a query builder for Metadata.
func (c *MetadataClient) Query() *MetadataQuery {
	return &MetadataQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeMetadata},
		inters: c.Interceptors(),
	}
}

// Get returns a Metadata entity by its id.
func (c *MetadataClient) Get(ctx context.Context, id int64) (*Metadata, error) {
	return c.Query().Where(metadata.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *MetadataClient) GetX(ctx context.Context, id int64) *Metadata {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRecord queries the record edge of a Metadata.
func (c *MetadataClient) QueryRecord(m *Metadata) *RecordQuery {
	query := (&RecordClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := m.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(metadata.Table, metadata.FieldID, id),
			sqlgraph.To(record.Table, record.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, metadata.RecordTable, metadata.RecordColumn),
		)
		fromV = sqlgraph.Neighbors(m.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryMetadataFormat queries the metadata_format edge of a Metadata.
func (c *MetadataClient) QueryMetadataFormat(m *Metadata) *MetadataFormatQuery {
	query := (&MetadataFormatClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := m.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(metadata.Table, metadata.FieldID, id),
			sqlgraph.To(metadataformat.Table, metadataformat.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, metadata.MetadataFormatTable, metadata.MetadataFormatColumn),
		)
		fromV = sqlgraph.Neighbors(m.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *MetadataClient) Hooks() []Hook {
	return c.hooks.Metadata
}

// Interceptors returns the client interceptors.
func (c *MetadataClient) Interceptors() []Interceptor {
	return c.inters.Metadata
}

func (c *MetadataClient) mutate(ctx context.Context, m *MetadataMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&MetadataCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&MetadataUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&MetadataUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&MetadataDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Metadata mutation op: %q", m.Op())
	}
}

// MetadataFormatClient is a client for the MetadataFormat schema.
type MetadataFormatClient struct {
	config
}

// NewMetadataFormatClient returns a client for the MetadataFormat from the given config.
func NewMetadataFormatClient(c config) *MetadataFormatClient {
	return &MetadataFormatClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `metadataformat.Hooks(f(g(h())))`.
func (c *MetadataFormatClient) Use(hooks ...Hook) {
	c.hooks.MetadataFormat = append(c.hooks.MetadataFormat, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `metadataformat.Intercept(f(g(h())))`.
func (c *MetadataFormatClient) Intercept(interceptors ...Interceptor) {
	c.inters.MetadataFormat = append(c.inters.MetadataFormat, interceptors...)
}

// Create returns a builder for creating a MetadataFormat entity.
func (c *MetadataFormatClient) Create() *MetadataFormatCreate {
	mutation := newMetadataFormatMutation(c.config, OpCreate)
	return &MetadataFormatCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of MetadataFormat entities.
func (c *MetadataFormatClient) CreateBulk(builders ...*MetadataFormatCreate) *MetadataFormatCreateBulk {
	return &MetadataFormatCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for MetadataFormat.
func (c *MetadataFormatClient) Update() *MetadataFormatUpdate {
	mutation := newMetadataFormatMutation(c.config, OpUpdate)
	return &MetadataFormatUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *MetadataFormatClient) UpdateOne(mf *MetadataFormat) *MetadataFormatUpdateOne {
	mutation := newMetadataFormatMutation(c.config, OpUpdateOne, withMetadataFormat(mf))
	return &MetadataFormatUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *MetadataFormatClient) UpdateOneID(id int64) *MetadataFormatUpdateOne {
	mutation := newMetadataFormatMutation(c.config, OpUpdateOne, withMetadataFormatID(id))
	return &MetadataFormatUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for MetadataFormat.
func (c *MetadataFormatClient) Delete() *MetadataFormatDelete {
	mutation := newMetadataFormatMutation(c.config, OpDelete)
	return &MetadataFormatDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *MetadataFormatClient) DeleteOne(mf *MetadataFormat) *MetadataFormatDeleteOne {
	return c.DeleteOneID(mf.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *MetadataFormatClient) DeleteOneID(id int64) *MetadataFormatDeleteOne {
	builder := c.Delete().Where(metadataformat.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &MetadataFormatDeleteOne{builder}
}

// Query returns a query builder for MetadataFormat.
func (c *MetadataFormatClient) Query() *MetadataFormatQuery {
	return &MetadataFormatQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeMetadataFormat},
		inters: c.Interceptors(),
	}
}

// Get returns a MetadataFormat entity by its id.
func (c *MetadataFormatClient) Get(ctx context.Context, id int64) (*MetadataFormat, error) {
	return c.Query().Where(metadataformat.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *MetadataFormatClient) GetX(ctx context.Context, id int64) *MetadataFormat {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryMetadata queries the metadata edge of a MetadataFormat.
func (c *MetadataFormatClient) QueryMetadata(mf *MetadataFormat) *MetadataQuery {
	query := (&MetadataClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := mf.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(metadataformat.Table, metadataformat.FieldID, id),
			sqlgraph.To(metadata.Table, metadata.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, metadataformat.MetadataTable, metadataformat.MetadataColumn),
		)
		fromV = sqlgraph.Neighbors(mf.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *MetadataFormatClient) Hooks() []Hook {
	return c.hooks.MetadataFormat
}

// Interceptors returns the client interceptors.
func (c *MetadataFormatClient) Interceptors() []Interceptor {
	return c.inters.MetadataFormat
}

func (c *MetadataFormatClient) mutate(ctx context.Context, m *MetadataFormatMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&MetadataFormatCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&MetadataFormatUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&MetadataFormatUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&MetadataFormatDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown MetadataFormat mutation op: %q", m.Op())
	}
}

// RecordClient is a client for the Record schema.
type RecordClient struct {
	config
}

// NewRecordClient returns a client for the Record from the given config.
func NewRecordClient(c config) *RecordClient {
	return &RecordClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `record.Hooks(f(g(h())))`.
func (c *RecordClient) Use(hooks ...Hook) {
	c.hooks.Record = append(c.hooks.Record, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `record.Intercept(f(g(h())))`.
func (c *RecordClient) Intercept(interceptors ...Interceptor) {
	c.inters.Record = append(c.inters.Record, interceptors...)
}

// Create returns a builder for creating a Record entity.
func (c *RecordClient) Create() *RecordCreate {
	mutation := newRecordMutation(c.config, OpCreate)
	return &RecordCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Record entities.
func (c *RecordClient) CreateBulk(builders ...*RecordCreate) *RecordCreateBulk {
	return &RecordCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Record.
func (c *RecordClient) Update() *RecordUpdate {
	mutation := newRecordMutation(c.config, OpUpdate)
	return &RecordUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *RecordClient) UpdateOne(r *Record) *RecordUpdateOne {
	mutation := newRecordMutation(c.config, OpUpdateOne, withRecord(r))
	return &RecordUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *RecordClient) UpdateOneID(id int64) *RecordUpdateOne {
	mutation := newRecordMutation(c.config, OpUpdateOne, withRecordID(id))
	return &RecordUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Record.
func (c *RecordClient) Delete() *RecordDelete {
	mutation := newRecordMutation(c.config, OpDelete)
	return &RecordDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *RecordClient) DeleteOne(r *Record) *RecordDeleteOne {
	return c.DeleteOneID(r.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *RecordClient) DeleteOneID(id int64) *RecordDeleteOne {
	builder := c.Delete().Where(record.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &RecordDeleteOne{builder}
}

// Query returns a query builder for Record.
func (c *RecordClient) Query() *RecordQuery {
	return &RecordQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeRecord},
		inters: c.Interceptors(),
	}
}

// Get returns a Record entity by its id.
func (c *RecordClient) Get(ctx context.Context, id int64) (*Record, error) {
	return c.Query().Where(record.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RecordClient) GetX(ctx context.Context, id int64) *Record {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryMetadata queries the metadata edge of a Record.
func (c *RecordClient) QueryMetadata(r *Record) *MetadataQuery {
	query := (&MetadataClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(record.Table, record.FieldID, id),
			sqlgraph.To(metadata.Table, metadata.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, record.MetadataTable, record.MetadataColumn),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QuerySets queries the sets edge of a Record.
func (c *RecordClient) QuerySets(r *Record) *SetQuery {
	query := (&SetClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(record.Table, record.FieldID, id),
			sqlgraph.To(set.Table, set.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, record.SetsTable, record.SetsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *RecordClient) Hooks() []Hook {
	return c.hooks.Record
}

// Interceptors returns the client interceptors.
func (c *RecordClient) Interceptors() []Interceptor {
	return c.inters.Record
}

func (c *RecordClient) mutate(ctx context.Context, m *RecordMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&RecordCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&RecordUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&RecordUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&RecordDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Record mutation op: %q", m.Op())
	}
}

// SetClient is a client for the Set schema.
type SetClient struct {
	config
}

// NewSetClient returns a client for the Set from the given config.
func NewSetClient(c config) *SetClient {
	return &SetClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `set.Hooks(f(g(h())))`.
func (c *SetClient) Use(hooks ...Hook) {
	c.hooks.Set = append(c.hooks.Set, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `set.Intercept(f(g(h())))`.
func (c *SetClient) Intercept(interceptors ...Interceptor) {
	c.inters.Set = append(c.inters.Set, interceptors...)
}

// Create returns a builder for creating a Set entity.
func (c *SetClient) Create() *SetCreate {
	mutation := newSetMutation(c.config, OpCreate)
	return &SetCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Set entities.
func (c *SetClient) CreateBulk(builders ...*SetCreate) *SetCreateBulk {
	return &SetCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Set.
func (c *SetClient) Update() *SetUpdate {
	mutation := newSetMutation(c.config, OpUpdate)
	return &SetUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SetClient) UpdateOne(s *Set) *SetUpdateOne {
	mutation := newSetMutation(c.config, OpUpdateOne, withSet(s))
	return &SetUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SetClient) UpdateOneID(id int64) *SetUpdateOne {
	mutation := newSetMutation(c.config, OpUpdateOne, withSetID(id))
	return &SetUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Set.
func (c *SetClient) Delete() *SetDelete {
	mutation := newSetMutation(c.config, OpDelete)
	return &SetDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SetClient) DeleteOne(s *Set) *SetDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SetClient) DeleteOneID(id int64) *SetDeleteOne {
	builder := c.Delete().Where(set.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SetDeleteOne{builder}
}

// Query returns a query builder for Set.
func (c *SetClient) Query() *SetQuery {
	return &SetQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSet},
		inters: c.Interceptors(),
	}
}

// Get returns a Set entity by its id.
func (c *SetClient) Get(ctx context.Context, id int64) (*Set, error) {
	return c.Query().Where(set.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SetClient) GetX(ctx context.Context, id int64) *Set {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRecords queries the records edge of a Set.
func (c *SetClient) QueryRecords(s *Set) *RecordQuery {
	query := (&RecordClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(set.Table, set.FieldID, id),
			sqlgraph.To(record.Table, record.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, set.RecordsTable, set.RecordsPrimaryKey...),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SetClient) Hooks() []Hook {
	return c.hooks.Set
}

// Interceptors returns the client interceptors.
func (c *SetClient) Interceptors() []Interceptor {
	return c.inters.Set
}

func (c *SetClient) mutate(ctx context.Context, m *SetMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SetCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SetUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SetUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SetDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Set mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Metadata, MetadataFormat, Record, Set []ent.Hook
	}
	inters struct {
		Metadata, MetadataFormat, Record, Set []ent.Interceptor
	}
)
