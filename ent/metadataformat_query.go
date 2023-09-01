// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/predicate"
	"github.com/ugent-library/oai-service/ent/record"
)

// MetadataFormatQuery is the builder for querying MetadataFormat entities.
type MetadataFormatQuery struct {
	config
	ctx         *QueryContext
	order       []metadataformat.OrderOption
	inters      []Interceptor
	predicates  []predicate.MetadataFormat
	withRecords *RecordQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the MetadataFormatQuery builder.
func (mfq *MetadataFormatQuery) Where(ps ...predicate.MetadataFormat) *MetadataFormatQuery {
	mfq.predicates = append(mfq.predicates, ps...)
	return mfq
}

// Limit the number of records to be returned by this query.
func (mfq *MetadataFormatQuery) Limit(limit int) *MetadataFormatQuery {
	mfq.ctx.Limit = &limit
	return mfq
}

// Offset to start from.
func (mfq *MetadataFormatQuery) Offset(offset int) *MetadataFormatQuery {
	mfq.ctx.Offset = &offset
	return mfq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mfq *MetadataFormatQuery) Unique(unique bool) *MetadataFormatQuery {
	mfq.ctx.Unique = &unique
	return mfq
}

// Order specifies how the records should be ordered.
func (mfq *MetadataFormatQuery) Order(o ...metadataformat.OrderOption) *MetadataFormatQuery {
	mfq.order = append(mfq.order, o...)
	return mfq
}

// QueryRecords chains the current query on the "records" edge.
func (mfq *MetadataFormatQuery) QueryRecords() *RecordQuery {
	query := (&RecordClient{config: mfq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mfq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mfq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(metadataformat.Table, metadataformat.FieldID, selector),
			sqlgraph.To(record.Table, record.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, metadataformat.RecordsTable, metadataformat.RecordsColumn),
		)
		fromU = sqlgraph.SetNeighbors(mfq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first MetadataFormat entity from the query.
// Returns a *NotFoundError when no MetadataFormat was found.
func (mfq *MetadataFormatQuery) First(ctx context.Context) (*MetadataFormat, error) {
	nodes, err := mfq.Limit(1).All(setContextOp(ctx, mfq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{metadataformat.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mfq *MetadataFormatQuery) FirstX(ctx context.Context) *MetadataFormat {
	node, err := mfq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first MetadataFormat ID from the query.
// Returns a *NotFoundError when no MetadataFormat ID was found.
func (mfq *MetadataFormatQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = mfq.Limit(1).IDs(setContextOp(ctx, mfq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{metadataformat.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (mfq *MetadataFormatQuery) FirstIDX(ctx context.Context) string {
	id, err := mfq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single MetadataFormat entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one MetadataFormat entity is found.
// Returns a *NotFoundError when no MetadataFormat entities are found.
func (mfq *MetadataFormatQuery) Only(ctx context.Context) (*MetadataFormat, error) {
	nodes, err := mfq.Limit(2).All(setContextOp(ctx, mfq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{metadataformat.Label}
	default:
		return nil, &NotSingularError{metadataformat.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mfq *MetadataFormatQuery) OnlyX(ctx context.Context) *MetadataFormat {
	node, err := mfq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only MetadataFormat ID in the query.
// Returns a *NotSingularError when more than one MetadataFormat ID is found.
// Returns a *NotFoundError when no entities are found.
func (mfq *MetadataFormatQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = mfq.Limit(2).IDs(setContextOp(ctx, mfq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{metadataformat.Label}
	default:
		err = &NotSingularError{metadataformat.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (mfq *MetadataFormatQuery) OnlyIDX(ctx context.Context) string {
	id, err := mfq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of MetadataFormats.
func (mfq *MetadataFormatQuery) All(ctx context.Context) ([]*MetadataFormat, error) {
	ctx = setContextOp(ctx, mfq.ctx, "All")
	if err := mfq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*MetadataFormat, *MetadataFormatQuery]()
	return withInterceptors[[]*MetadataFormat](ctx, mfq, qr, mfq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mfq *MetadataFormatQuery) AllX(ctx context.Context) []*MetadataFormat {
	nodes, err := mfq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of MetadataFormat IDs.
func (mfq *MetadataFormatQuery) IDs(ctx context.Context) (ids []string, err error) {
	if mfq.ctx.Unique == nil && mfq.path != nil {
		mfq.Unique(true)
	}
	ctx = setContextOp(ctx, mfq.ctx, "IDs")
	if err = mfq.Select(metadataformat.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (mfq *MetadataFormatQuery) IDsX(ctx context.Context) []string {
	ids, err := mfq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (mfq *MetadataFormatQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mfq.ctx, "Count")
	if err := mfq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mfq, querierCount[*MetadataFormatQuery](), mfq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mfq *MetadataFormatQuery) CountX(ctx context.Context) int {
	count, err := mfq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mfq *MetadataFormatQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mfq.ctx, "Exist")
	switch _, err := mfq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mfq *MetadataFormatQuery) ExistX(ctx context.Context) bool {
	exist, err := mfq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the MetadataFormatQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mfq *MetadataFormatQuery) Clone() *MetadataFormatQuery {
	if mfq == nil {
		return nil
	}
	return &MetadataFormatQuery{
		config:      mfq.config,
		ctx:         mfq.ctx.Clone(),
		order:       append([]metadataformat.OrderOption{}, mfq.order...),
		inters:      append([]Interceptor{}, mfq.inters...),
		predicates:  append([]predicate.MetadataFormat{}, mfq.predicates...),
		withRecords: mfq.withRecords.Clone(),
		// clone intermediate query.
		sql:  mfq.sql.Clone(),
		path: mfq.path,
	}
}

// WithRecords tells the query-builder to eager-load the nodes that are connected to
// the "records" edge. The optional arguments are used to configure the query builder of the edge.
func (mfq *MetadataFormatQuery) WithRecords(opts ...func(*RecordQuery)) *MetadataFormatQuery {
	query := (&RecordClient{config: mfq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mfq.withRecords = query
	return mfq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Schema string `json:"schema,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.MetadataFormat.Query().
//		GroupBy(metadataformat.FieldSchema).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (mfq *MetadataFormatQuery) GroupBy(field string, fields ...string) *MetadataFormatGroupBy {
	mfq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &MetadataFormatGroupBy{build: mfq}
	grbuild.flds = &mfq.ctx.Fields
	grbuild.label = metadataformat.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Schema string `json:"schema,omitempty"`
//	}
//
//	client.MetadataFormat.Query().
//		Select(metadataformat.FieldSchema).
//		Scan(ctx, &v)
func (mfq *MetadataFormatQuery) Select(fields ...string) *MetadataFormatSelect {
	mfq.ctx.Fields = append(mfq.ctx.Fields, fields...)
	sbuild := &MetadataFormatSelect{MetadataFormatQuery: mfq}
	sbuild.label = metadataformat.Label
	sbuild.flds, sbuild.scan = &mfq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a MetadataFormatSelect configured with the given aggregations.
func (mfq *MetadataFormatQuery) Aggregate(fns ...AggregateFunc) *MetadataFormatSelect {
	return mfq.Select().Aggregate(fns...)
}

func (mfq *MetadataFormatQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mfq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mfq); err != nil {
				return err
			}
		}
	}
	for _, f := range mfq.ctx.Fields {
		if !metadataformat.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mfq.path != nil {
		prev, err := mfq.path(ctx)
		if err != nil {
			return err
		}
		mfq.sql = prev
	}
	return nil
}

func (mfq *MetadataFormatQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*MetadataFormat, error) {
	var (
		nodes       = []*MetadataFormat{}
		_spec       = mfq.querySpec()
		loadedTypes = [1]bool{
			mfq.withRecords != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*MetadataFormat).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &MetadataFormat{config: mfq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mfq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := mfq.withRecords; query != nil {
		if err := mfq.loadRecords(ctx, query, nodes,
			func(n *MetadataFormat) { n.Edges.Records = []*Record{} },
			func(n *MetadataFormat, e *Record) { n.Edges.Records = append(n.Edges.Records, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (mfq *MetadataFormatQuery) loadRecords(ctx context.Context, query *RecordQuery, nodes []*MetadataFormat, init func(*MetadataFormat), assign func(*MetadataFormat, *Record)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[string]*MetadataFormat)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(record.FieldMetadataFormatID)
	}
	query.Where(predicate.Record(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(metadataformat.RecordsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.MetadataFormatID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "metadata_format_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (mfq *MetadataFormatQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mfq.querySpec()
	_spec.Node.Columns = mfq.ctx.Fields
	if len(mfq.ctx.Fields) > 0 {
		_spec.Unique = mfq.ctx.Unique != nil && *mfq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, mfq.driver, _spec)
}

func (mfq *MetadataFormatQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(metadataformat.Table, metadataformat.Columns, sqlgraph.NewFieldSpec(metadataformat.FieldID, field.TypeString))
	_spec.From = mfq.sql
	if unique := mfq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if mfq.path != nil {
		_spec.Unique = true
	}
	if fields := mfq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, metadataformat.FieldID)
		for i := range fields {
			if fields[i] != metadataformat.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := mfq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mfq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mfq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mfq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mfq *MetadataFormatQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mfq.driver.Dialect())
	t1 := builder.Table(metadataformat.Table)
	columns := mfq.ctx.Fields
	if len(columns) == 0 {
		columns = metadataformat.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mfq.sql != nil {
		selector = mfq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mfq.ctx.Unique != nil && *mfq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range mfq.predicates {
		p(selector)
	}
	for _, p := range mfq.order {
		p(selector)
	}
	if offset := mfq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mfq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// MetadataFormatGroupBy is the group-by builder for MetadataFormat entities.
type MetadataFormatGroupBy struct {
	selector
	build *MetadataFormatQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mfgb *MetadataFormatGroupBy) Aggregate(fns ...AggregateFunc) *MetadataFormatGroupBy {
	mfgb.fns = append(mfgb.fns, fns...)
	return mfgb
}

// Scan applies the selector query and scans the result into the given value.
func (mfgb *MetadataFormatGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mfgb.build.ctx, "GroupBy")
	if err := mfgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MetadataFormatQuery, *MetadataFormatGroupBy](ctx, mfgb.build, mfgb, mfgb.build.inters, v)
}

func (mfgb *MetadataFormatGroupBy) sqlScan(ctx context.Context, root *MetadataFormatQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mfgb.fns))
	for _, fn := range mfgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mfgb.flds)+len(mfgb.fns))
		for _, f := range *mfgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mfgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mfgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// MetadataFormatSelect is the builder for selecting fields of MetadataFormat entities.
type MetadataFormatSelect struct {
	*MetadataFormatQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (mfs *MetadataFormatSelect) Aggregate(fns ...AggregateFunc) *MetadataFormatSelect {
	mfs.fns = append(mfs.fns, fns...)
	return mfs
}

// Scan applies the selector query and scans the result into the given value.
func (mfs *MetadataFormatSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mfs.ctx, "Select")
	if err := mfs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*MetadataFormatQuery, *MetadataFormatSelect](ctx, mfs.MetadataFormatQuery, mfs, mfs.inters, v)
}

func (mfs *MetadataFormatSelect) sqlScan(ctx context.Context, root *MetadataFormatQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(mfs.fns))
	for _, fn := range mfs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*mfs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mfs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
