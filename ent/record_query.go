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
	"github.com/ugent-library/oai-service/ent/set"
)

// RecordQuery is the builder for querying Record entities.
type RecordQuery struct {
	config
	limit              *int
	offset             *int
	unique             *bool
	order              []OrderFunc
	fields             []string
	predicates         []predicate.Record
	withMetadataFormat *MetadataFormatQuery
	withSets           *SetQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the RecordQuery builder.
func (rq *RecordQuery) Where(ps ...predicate.Record) *RecordQuery {
	rq.predicates = append(rq.predicates, ps...)
	return rq
}

// Limit adds a limit step to the query.
func (rq *RecordQuery) Limit(limit int) *RecordQuery {
	rq.limit = &limit
	return rq
}

// Offset adds an offset step to the query.
func (rq *RecordQuery) Offset(offset int) *RecordQuery {
	rq.offset = &offset
	return rq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (rq *RecordQuery) Unique(unique bool) *RecordQuery {
	rq.unique = &unique
	return rq
}

// Order adds an order step to the query.
func (rq *RecordQuery) Order(o ...OrderFunc) *RecordQuery {
	rq.order = append(rq.order, o...)
	return rq
}

// QueryMetadataFormat chains the current query on the "metadata_format" edge.
func (rq *RecordQuery) QueryMetadataFormat() *MetadataFormatQuery {
	query := &MetadataFormatQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(record.Table, record.FieldID, selector),
			sqlgraph.To(metadataformat.Table, metadataformat.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, record.MetadataFormatTable, record.MetadataFormatColumn),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySets chains the current query on the "sets" edge.
func (rq *RecordQuery) QuerySets() *SetQuery {
	query := &SetQuery{config: rq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := rq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(record.Table, record.FieldID, selector),
			sqlgraph.To(set.Table, set.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, record.SetsTable, record.SetsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(rq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Record entity from the query.
// Returns a *NotFoundError when no Record was found.
func (rq *RecordQuery) First(ctx context.Context) (*Record, error) {
	nodes, err := rq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{record.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (rq *RecordQuery) FirstX(ctx context.Context) *Record {
	node, err := rq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Record ID from the query.
// Returns a *NotFoundError when no Record ID was found.
func (rq *RecordQuery) FirstID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = rq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{record.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (rq *RecordQuery) FirstIDX(ctx context.Context) int64 {
	id, err := rq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Record entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Record entity is found.
// Returns a *NotFoundError when no Record entities are found.
func (rq *RecordQuery) Only(ctx context.Context) (*Record, error) {
	nodes, err := rq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{record.Label}
	default:
		return nil, &NotSingularError{record.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (rq *RecordQuery) OnlyX(ctx context.Context) *Record {
	node, err := rq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Record ID in the query.
// Returns a *NotSingularError when more than one Record ID is found.
// Returns a *NotFoundError when no entities are found.
func (rq *RecordQuery) OnlyID(ctx context.Context) (id int64, err error) {
	var ids []int64
	if ids, err = rq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{record.Label}
	default:
		err = &NotSingularError{record.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (rq *RecordQuery) OnlyIDX(ctx context.Context) int64 {
	id, err := rq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Records.
func (rq *RecordQuery) All(ctx context.Context) ([]*Record, error) {
	if err := rq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return rq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (rq *RecordQuery) AllX(ctx context.Context) []*Record {
	nodes, err := rq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Record IDs.
func (rq *RecordQuery) IDs(ctx context.Context) ([]int64, error) {
	var ids []int64
	if err := rq.Select(record.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (rq *RecordQuery) IDsX(ctx context.Context) []int64 {
	ids, err := rq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (rq *RecordQuery) Count(ctx context.Context) (int, error) {
	if err := rq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return rq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (rq *RecordQuery) CountX(ctx context.Context) int {
	count, err := rq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (rq *RecordQuery) Exist(ctx context.Context) (bool, error) {
	if err := rq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return rq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (rq *RecordQuery) ExistX(ctx context.Context) bool {
	exist, err := rq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the RecordQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (rq *RecordQuery) Clone() *RecordQuery {
	if rq == nil {
		return nil
	}
	return &RecordQuery{
		config:             rq.config,
		limit:              rq.limit,
		offset:             rq.offset,
		order:              append([]OrderFunc{}, rq.order...),
		predicates:         append([]predicate.Record{}, rq.predicates...),
		withMetadataFormat: rq.withMetadataFormat.Clone(),
		withSets:           rq.withSets.Clone(),
		// clone intermediate query.
		sql:    rq.sql.Clone(),
		path:   rq.path,
		unique: rq.unique,
	}
}

// WithMetadataFormat tells the query-builder to eager-load the nodes that are connected to
// the "metadata_format" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RecordQuery) WithMetadataFormat(opts ...func(*MetadataFormatQuery)) *RecordQuery {
	query := &MetadataFormatQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withMetadataFormat = query
	return rq
}

// WithSets tells the query-builder to eager-load the nodes that are connected to
// the "sets" edge. The optional arguments are used to configure the query builder of the edge.
func (rq *RecordQuery) WithSets(opts ...func(*SetQuery)) *RecordQuery {
	query := &SetQuery{config: rq.config}
	for _, opt := range opts {
		opt(query)
	}
	rq.withSets = query
	return rq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		MetadataFormatID int64 `json:"metadata_format_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Record.Query().
//		GroupBy(record.FieldMetadataFormatID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (rq *RecordQuery) GroupBy(field string, fields ...string) *RecordGroupBy {
	grbuild := &RecordGroupBy{config: rq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := rq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return rq.sqlQuery(ctx), nil
	}
	grbuild.label = record.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		MetadataFormatID int64 `json:"metadata_format_id,omitempty"`
//	}
//
//	client.Record.Query().
//		Select(record.FieldMetadataFormatID).
//		Scan(ctx, &v)
func (rq *RecordQuery) Select(fields ...string) *RecordSelect {
	rq.fields = append(rq.fields, fields...)
	selbuild := &RecordSelect{RecordQuery: rq}
	selbuild.label = record.Label
	selbuild.flds, selbuild.scan = &rq.fields, selbuild.Scan
	return selbuild
}

// Aggregate returns a RecordSelect configured with the given aggregations.
func (rq *RecordQuery) Aggregate(fns ...AggregateFunc) *RecordSelect {
	return rq.Select().Aggregate(fns...)
}

func (rq *RecordQuery) prepareQuery(ctx context.Context) error {
	for _, f := range rq.fields {
		if !record.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if rq.path != nil {
		prev, err := rq.path(ctx)
		if err != nil {
			return err
		}
		rq.sql = prev
	}
	return nil
}

func (rq *RecordQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Record, error) {
	var (
		nodes       = []*Record{}
		_spec       = rq.querySpec()
		loadedTypes = [2]bool{
			rq.withMetadataFormat != nil,
			rq.withSets != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Record).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Record{config: rq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, rq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := rq.withMetadataFormat; query != nil {
		if err := rq.loadMetadataFormat(ctx, query, nodes, nil,
			func(n *Record, e *MetadataFormat) { n.Edges.MetadataFormat = e }); err != nil {
			return nil, err
		}
	}
	if query := rq.withSets; query != nil {
		if err := rq.loadSets(ctx, query, nodes,
			func(n *Record) { n.Edges.Sets = []*Set{} },
			func(n *Record, e *Set) { n.Edges.Sets = append(n.Edges.Sets, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (rq *RecordQuery) loadMetadataFormat(ctx context.Context, query *MetadataFormatQuery, nodes []*Record, init func(*Record), assign func(*Record, *MetadataFormat)) error {
	ids := make([]int64, 0, len(nodes))
	nodeids := make(map[int64][]*Record)
	for i := range nodes {
		fk := nodes[i].MetadataFormatID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(metadataformat.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "metadata_format_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (rq *RecordQuery) loadSets(ctx context.Context, query *SetQuery, nodes []*Record, init func(*Record), assign func(*Record, *Set)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int64]*Record)
	nids := make(map[int64]map[*Record]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(record.SetsTable)
		s.Join(joinT).On(s.C(set.FieldID), joinT.C(record.SetsPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(record.SetsPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(record.SetsPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	neighbors, err := query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
		assign := spec.Assign
		values := spec.ScanValues
		spec.ScanValues = func(columns []string) ([]any, error) {
			values, err := values(columns[1:])
			if err != nil {
				return nil, err
			}
			return append([]any{new(sql.NullInt64)}, values...), nil
		}
		spec.Assign = func(columns []string, values []any) error {
			outValue := values[0].(*sql.NullInt64).Int64
			inValue := values[1].(*sql.NullInt64).Int64
			if nids[inValue] == nil {
				nids[inValue] = map[*Record]struct{}{byID[outValue]: {}}
				return assign(columns[1:], values[1:])
			}
			nids[inValue][byID[outValue]] = struct{}{}
			return nil
		}
	})
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "sets" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (rq *RecordQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := rq.querySpec()
	_spec.Node.Columns = rq.fields
	if len(rq.fields) > 0 {
		_spec.Unique = rq.unique != nil && *rq.unique
	}
	return sqlgraph.CountNodes(ctx, rq.driver, _spec)
}

func (rq *RecordQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := rq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (rq *RecordQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   record.Table,
			Columns: record.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt64,
				Column: record.FieldID,
			},
		},
		From:   rq.sql,
		Unique: true,
	}
	if unique := rq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := rq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, record.FieldID)
		for i := range fields {
			if fields[i] != record.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := rq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := rq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := rq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := rq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (rq *RecordQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(rq.driver.Dialect())
	t1 := builder.Table(record.Table)
	columns := rq.fields
	if len(columns) == 0 {
		columns = record.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if rq.sql != nil {
		selector = rq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if rq.unique != nil && *rq.unique {
		selector.Distinct()
	}
	for _, p := range rq.predicates {
		p(selector)
	}
	for _, p := range rq.order {
		p(selector)
	}
	if offset := rq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := rq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// RecordGroupBy is the group-by builder for Record entities.
type RecordGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (rgb *RecordGroupBy) Aggregate(fns ...AggregateFunc) *RecordGroupBy {
	rgb.fns = append(rgb.fns, fns...)
	return rgb
}

// Scan applies the group-by query and scans the result into the given value.
func (rgb *RecordGroupBy) Scan(ctx context.Context, v any) error {
	query, err := rgb.path(ctx)
	if err != nil {
		return err
	}
	rgb.sql = query
	return rgb.sqlScan(ctx, v)
}

func (rgb *RecordGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range rgb.fields {
		if !record.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := rgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := rgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (rgb *RecordGroupBy) sqlQuery() *sql.Selector {
	selector := rgb.sql.Select()
	aggregation := make([]string, 0, len(rgb.fns))
	for _, fn := range rgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(rgb.fields)+len(rgb.fns))
		for _, f := range rgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(rgb.fields...)...)
}

// RecordSelect is the builder for selecting fields of Record entities.
type RecordSelect struct {
	*RecordQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (rs *RecordSelect) Aggregate(fns ...AggregateFunc) *RecordSelect {
	rs.fns = append(rs.fns, fns...)
	return rs
}

// Scan applies the selector query and scans the result into the given value.
func (rs *RecordSelect) Scan(ctx context.Context, v any) error {
	if err := rs.prepareQuery(ctx); err != nil {
		return err
	}
	rs.sql = rs.RecordQuery.sqlQuery(ctx)
	return rs.sqlScan(ctx, v)
}

func (rs *RecordSelect) sqlScan(ctx context.Context, v any) error {
	aggregation := make([]string, 0, len(rs.fns))
	for _, fn := range rs.fns {
		aggregation = append(aggregation, fn(rs.sql))
	}
	switch n := len(*rs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		rs.sql.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		rs.sql.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := rs.sql.Query()
	if err := rs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
