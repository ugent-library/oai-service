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
	"github.com/ugent-library/oai-service/ent/item"
	"github.com/ugent-library/oai-service/ent/predicate"
	"github.com/ugent-library/oai-service/ent/record"
	"github.com/ugent-library/oai-service/ent/set"
)

// ItemQuery is the builder for querying Item entities.
type ItemQuery struct {
	config
	ctx         *QueryContext
	order       []item.OrderOption
	inters      []Interceptor
	predicates  []predicate.Item
	withRecords *RecordQuery
	withSets    *SetQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ItemQuery builder.
func (iq *ItemQuery) Where(ps ...predicate.Item) *ItemQuery {
	iq.predicates = append(iq.predicates, ps...)
	return iq
}

// Limit the number of records to be returned by this query.
func (iq *ItemQuery) Limit(limit int) *ItemQuery {
	iq.ctx.Limit = &limit
	return iq
}

// Offset to start from.
func (iq *ItemQuery) Offset(offset int) *ItemQuery {
	iq.ctx.Offset = &offset
	return iq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (iq *ItemQuery) Unique(unique bool) *ItemQuery {
	iq.ctx.Unique = &unique
	return iq
}

// Order specifies how the records should be ordered.
func (iq *ItemQuery) Order(o ...item.OrderOption) *ItemQuery {
	iq.order = append(iq.order, o...)
	return iq
}

// QueryRecords chains the current query on the "records" edge.
func (iq *ItemQuery) QueryRecords() *RecordQuery {
	query := (&RecordClient{config: iq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := iq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := iq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(item.Table, item.FieldID, selector),
			sqlgraph.To(record.Table, record.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, item.RecordsTable, item.RecordsColumn),
		)
		fromU = sqlgraph.SetNeighbors(iq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QuerySets chains the current query on the "sets" edge.
func (iq *ItemQuery) QuerySets() *SetQuery {
	query := (&SetClient{config: iq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := iq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := iq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(item.Table, item.FieldID, selector),
			sqlgraph.To(set.Table, set.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, item.SetsTable, item.SetsPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(iq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Item entity from the query.
// Returns a *NotFoundError when no Item was found.
func (iq *ItemQuery) First(ctx context.Context) (*Item, error) {
	nodes, err := iq.Limit(1).All(setContextOp(ctx, iq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{item.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (iq *ItemQuery) FirstX(ctx context.Context) *Item {
	node, err := iq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Item ID from the query.
// Returns a *NotFoundError when no Item ID was found.
func (iq *ItemQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = iq.Limit(1).IDs(setContextOp(ctx, iq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{item.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (iq *ItemQuery) FirstIDX(ctx context.Context) string {
	id, err := iq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Item entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Item entity is found.
// Returns a *NotFoundError when no Item entities are found.
func (iq *ItemQuery) Only(ctx context.Context) (*Item, error) {
	nodes, err := iq.Limit(2).All(setContextOp(ctx, iq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{item.Label}
	default:
		return nil, &NotSingularError{item.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (iq *ItemQuery) OnlyX(ctx context.Context) *Item {
	node, err := iq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Item ID in the query.
// Returns a *NotSingularError when more than one Item ID is found.
// Returns a *NotFoundError when no entities are found.
func (iq *ItemQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = iq.Limit(2).IDs(setContextOp(ctx, iq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{item.Label}
	default:
		err = &NotSingularError{item.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (iq *ItemQuery) OnlyIDX(ctx context.Context) string {
	id, err := iq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Items.
func (iq *ItemQuery) All(ctx context.Context) ([]*Item, error) {
	ctx = setContextOp(ctx, iq.ctx, "All")
	if err := iq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Item, *ItemQuery]()
	return withInterceptors[[]*Item](ctx, iq, qr, iq.inters)
}

// AllX is like All, but panics if an error occurs.
func (iq *ItemQuery) AllX(ctx context.Context) []*Item {
	nodes, err := iq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Item IDs.
func (iq *ItemQuery) IDs(ctx context.Context) (ids []string, err error) {
	if iq.ctx.Unique == nil && iq.path != nil {
		iq.Unique(true)
	}
	ctx = setContextOp(ctx, iq.ctx, "IDs")
	if err = iq.Select(item.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (iq *ItemQuery) IDsX(ctx context.Context) []string {
	ids, err := iq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (iq *ItemQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, iq.ctx, "Count")
	if err := iq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, iq, querierCount[*ItemQuery](), iq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (iq *ItemQuery) CountX(ctx context.Context) int {
	count, err := iq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (iq *ItemQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, iq.ctx, "Exist")
	switch _, err := iq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (iq *ItemQuery) ExistX(ctx context.Context) bool {
	exist, err := iq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ItemQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (iq *ItemQuery) Clone() *ItemQuery {
	if iq == nil {
		return nil
	}
	return &ItemQuery{
		config:      iq.config,
		ctx:         iq.ctx.Clone(),
		order:       append([]item.OrderOption{}, iq.order...),
		inters:      append([]Interceptor{}, iq.inters...),
		predicates:  append([]predicate.Item{}, iq.predicates...),
		withRecords: iq.withRecords.Clone(),
		withSets:    iq.withSets.Clone(),
		// clone intermediate query.
		sql:  iq.sql.Clone(),
		path: iq.path,
	}
}

// WithRecords tells the query-builder to eager-load the nodes that are connected to
// the "records" edge. The optional arguments are used to configure the query builder of the edge.
func (iq *ItemQuery) WithRecords(opts ...func(*RecordQuery)) *ItemQuery {
	query := (&RecordClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	iq.withRecords = query
	return iq
}

// WithSets tells the query-builder to eager-load the nodes that are connected to
// the "sets" edge. The optional arguments are used to configure the query builder of the edge.
func (iq *ItemQuery) WithSets(opts ...func(*SetQuery)) *ItemQuery {
	query := (&SetClient{config: iq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	iq.withSets = query
	return iq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
func (iq *ItemQuery) GroupBy(field string, fields ...string) *ItemGroupBy {
	iq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ItemGroupBy{build: iq}
	grbuild.flds = &iq.ctx.Fields
	grbuild.label = item.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
func (iq *ItemQuery) Select(fields ...string) *ItemSelect {
	iq.ctx.Fields = append(iq.ctx.Fields, fields...)
	sbuild := &ItemSelect{ItemQuery: iq}
	sbuild.label = item.Label
	sbuild.flds, sbuild.scan = &iq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ItemSelect configured with the given aggregations.
func (iq *ItemQuery) Aggregate(fns ...AggregateFunc) *ItemSelect {
	return iq.Select().Aggregate(fns...)
}

func (iq *ItemQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range iq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, iq); err != nil {
				return err
			}
		}
	}
	for _, f := range iq.ctx.Fields {
		if !item.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if iq.path != nil {
		prev, err := iq.path(ctx)
		if err != nil {
			return err
		}
		iq.sql = prev
	}
	return nil
}

func (iq *ItemQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Item, error) {
	var (
		nodes       = []*Item{}
		_spec       = iq.querySpec()
		loadedTypes = [2]bool{
			iq.withRecords != nil,
			iq.withSets != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Item).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Item{config: iq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, iq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := iq.withRecords; query != nil {
		if err := iq.loadRecords(ctx, query, nodes,
			func(n *Item) { n.Edges.Records = []*Record{} },
			func(n *Item, e *Record) { n.Edges.Records = append(n.Edges.Records, e) }); err != nil {
			return nil, err
		}
	}
	if query := iq.withSets; query != nil {
		if err := iq.loadSets(ctx, query, nodes,
			func(n *Item) { n.Edges.Sets = []*Set{} },
			func(n *Item, e *Set) { n.Edges.Sets = append(n.Edges.Sets, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (iq *ItemQuery) loadRecords(ctx context.Context, query *RecordQuery, nodes []*Item, init func(*Item), assign func(*Item, *Record)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[string]*Item)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(record.FieldItemID)
	}
	query.Where(predicate.Record(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(item.RecordsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.ItemID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "item_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (iq *ItemQuery) loadSets(ctx context.Context, query *SetQuery, nodes []*Item, init func(*Item), assign func(*Item, *Set)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[string]*Item)
	nids := make(map[string]map[*Item]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(item.SetsTable)
		s.Join(joinT).On(s.C(set.FieldID), joinT.C(item.SetsPrimaryKey[1]))
		s.Where(sql.InValues(joinT.C(item.SetsPrimaryKey[0]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(item.SetsPrimaryKey[0]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullString)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := values[0].(*sql.NullString).String
				inValue := values[1].(*sql.NullString).String
				if nids[inValue] == nil {
					nids[inValue] = map[*Item]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Set](ctx, query, qr, query.inters)
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

func (iq *ItemQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := iq.querySpec()
	_spec.Node.Columns = iq.ctx.Fields
	if len(iq.ctx.Fields) > 0 {
		_spec.Unique = iq.ctx.Unique != nil && *iq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, iq.driver, _spec)
}

func (iq *ItemQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(item.Table, item.Columns, sqlgraph.NewFieldSpec(item.FieldID, field.TypeString))
	_spec.From = iq.sql
	if unique := iq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if iq.path != nil {
		_spec.Unique = true
	}
	if fields := iq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, item.FieldID)
		for i := range fields {
			if fields[i] != item.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := iq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := iq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := iq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := iq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (iq *ItemQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(iq.driver.Dialect())
	t1 := builder.Table(item.Table)
	columns := iq.ctx.Fields
	if len(columns) == 0 {
		columns = item.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if iq.sql != nil {
		selector = iq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if iq.ctx.Unique != nil && *iq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range iq.predicates {
		p(selector)
	}
	for _, p := range iq.order {
		p(selector)
	}
	if offset := iq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := iq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ItemGroupBy is the group-by builder for Item entities.
type ItemGroupBy struct {
	selector
	build *ItemQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (igb *ItemGroupBy) Aggregate(fns ...AggregateFunc) *ItemGroupBy {
	igb.fns = append(igb.fns, fns...)
	return igb
}

// Scan applies the selector query and scans the result into the given value.
func (igb *ItemGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, igb.build.ctx, "GroupBy")
	if err := igb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ItemQuery, *ItemGroupBy](ctx, igb.build, igb, igb.build.inters, v)
}

func (igb *ItemGroupBy) sqlScan(ctx context.Context, root *ItemQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(igb.fns))
	for _, fn := range igb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*igb.flds)+len(igb.fns))
		for _, f := range *igb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*igb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := igb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ItemSelect is the builder for selecting fields of Item entities.
type ItemSelect struct {
	*ItemQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (is *ItemSelect) Aggregate(fns ...AggregateFunc) *ItemSelect {
	is.fns = append(is.fns, fns...)
	return is
}

// Scan applies the selector query and scans the result into the given value.
func (is *ItemSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, is.ctx, "Select")
	if err := is.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ItemQuery, *ItemSelect](ctx, is.ItemQuery, is, is.inters, v)
}

func (is *ItemSelect) sqlScan(ctx context.Context, root *ItemQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(is.fns))
	for _, fn := range is.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*is.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := is.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}