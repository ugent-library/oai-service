// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/record"
)

// MetadataFormatCreate is the builder for creating a MetadataFormat entity.
type MetadataFormatCreate struct {
	config
	mutation *MetadataFormatMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetSchema sets the "schema" field.
func (mfc *MetadataFormatCreate) SetSchema(s string) *MetadataFormatCreate {
	mfc.mutation.SetSchema(s)
	return mfc
}

// SetNamespace sets the "namespace" field.
func (mfc *MetadataFormatCreate) SetNamespace(s string) *MetadataFormatCreate {
	mfc.mutation.SetNamespace(s)
	return mfc
}

// SetID sets the "id" field.
func (mfc *MetadataFormatCreate) SetID(s string) *MetadataFormatCreate {
	mfc.mutation.SetID(s)
	return mfc
}

// AddRecordIDs adds the "records" edge to the Record entity by IDs.
func (mfc *MetadataFormatCreate) AddRecordIDs(ids ...int64) *MetadataFormatCreate {
	mfc.mutation.AddRecordIDs(ids...)
	return mfc
}

// AddRecords adds the "records" edges to the Record entity.
func (mfc *MetadataFormatCreate) AddRecords(r ...*Record) *MetadataFormatCreate {
	ids := make([]int64, len(r))
	for i := range r {
		ids[i] = r[i].ID
	}
	return mfc.AddRecordIDs(ids...)
}

// Mutation returns the MetadataFormatMutation object of the builder.
func (mfc *MetadataFormatCreate) Mutation() *MetadataFormatMutation {
	return mfc.mutation
}

// Save creates the MetadataFormat in the database.
func (mfc *MetadataFormatCreate) Save(ctx context.Context) (*MetadataFormat, error) {
	return withHooks(ctx, mfc.sqlSave, mfc.mutation, mfc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mfc *MetadataFormatCreate) SaveX(ctx context.Context) *MetadataFormat {
	v, err := mfc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mfc *MetadataFormatCreate) Exec(ctx context.Context) error {
	_, err := mfc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mfc *MetadataFormatCreate) ExecX(ctx context.Context) {
	if err := mfc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mfc *MetadataFormatCreate) check() error {
	if _, ok := mfc.mutation.Schema(); !ok {
		return &ValidationError{Name: "schema", err: errors.New(`ent: missing required field "MetadataFormat.schema"`)}
	}
	if _, ok := mfc.mutation.Namespace(); !ok {
		return &ValidationError{Name: "namespace", err: errors.New(`ent: missing required field "MetadataFormat.namespace"`)}
	}
	return nil
}

func (mfc *MetadataFormatCreate) sqlSave(ctx context.Context) (*MetadataFormat, error) {
	if err := mfc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mfc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mfc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected MetadataFormat.ID type: %T", _spec.ID.Value)
		}
	}
	mfc.mutation.id = &_node.ID
	mfc.mutation.done = true
	return _node, nil
}

func (mfc *MetadataFormatCreate) createSpec() (*MetadataFormat, *sqlgraph.CreateSpec) {
	var (
		_node = &MetadataFormat{config: mfc.config}
		_spec = sqlgraph.NewCreateSpec(metadataformat.Table, sqlgraph.NewFieldSpec(metadataformat.FieldID, field.TypeString))
	)
	_spec.OnConflict = mfc.conflict
	if id, ok := mfc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := mfc.mutation.Schema(); ok {
		_spec.SetField(metadataformat.FieldSchema, field.TypeString, value)
		_node.Schema = value
	}
	if value, ok := mfc.mutation.Namespace(); ok {
		_spec.SetField(metadataformat.FieldNamespace, field.TypeString, value)
		_node.Namespace = value
	}
	if nodes := mfc.mutation.RecordsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   metadataformat.RecordsTable,
			Columns: []string{metadataformat.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(record.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.MetadataFormat.Create().
//		SetSchema(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MetadataFormatUpsert) {
//			SetSchema(v+v).
//		}).
//		Exec(ctx)
func (mfc *MetadataFormatCreate) OnConflict(opts ...sql.ConflictOption) *MetadataFormatUpsertOne {
	mfc.conflict = opts
	return &MetadataFormatUpsertOne{
		create: mfc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.MetadataFormat.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mfc *MetadataFormatCreate) OnConflictColumns(columns ...string) *MetadataFormatUpsertOne {
	mfc.conflict = append(mfc.conflict, sql.ConflictColumns(columns...))
	return &MetadataFormatUpsertOne{
		create: mfc,
	}
}

type (
	// MetadataFormatUpsertOne is the builder for "upsert"-ing
	//  one MetadataFormat node.
	MetadataFormatUpsertOne struct {
		create *MetadataFormatCreate
	}

	// MetadataFormatUpsert is the "OnConflict" setter.
	MetadataFormatUpsert struct {
		*sql.UpdateSet
	}
)

// SetSchema sets the "schema" field.
func (u *MetadataFormatUpsert) SetSchema(v string) *MetadataFormatUpsert {
	u.Set(metadataformat.FieldSchema, v)
	return u
}

// UpdateSchema sets the "schema" field to the value that was provided on create.
func (u *MetadataFormatUpsert) UpdateSchema() *MetadataFormatUpsert {
	u.SetExcluded(metadataformat.FieldSchema)
	return u
}

// SetNamespace sets the "namespace" field.
func (u *MetadataFormatUpsert) SetNamespace(v string) *MetadataFormatUpsert {
	u.Set(metadataformat.FieldNamespace, v)
	return u
}

// UpdateNamespace sets the "namespace" field to the value that was provided on create.
func (u *MetadataFormatUpsert) UpdateNamespace() *MetadataFormatUpsert {
	u.SetExcluded(metadataformat.FieldNamespace)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.MetadataFormat.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(metadataformat.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MetadataFormatUpsertOne) UpdateNewValues() *MetadataFormatUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(metadataformat.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.MetadataFormat.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *MetadataFormatUpsertOne) Ignore() *MetadataFormatUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MetadataFormatUpsertOne) DoNothing() *MetadataFormatUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MetadataFormatCreate.OnConflict
// documentation for more info.
func (u *MetadataFormatUpsertOne) Update(set func(*MetadataFormatUpsert)) *MetadataFormatUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MetadataFormatUpsert{UpdateSet: update})
	}))
	return u
}

// SetSchema sets the "schema" field.
func (u *MetadataFormatUpsertOne) SetSchema(v string) *MetadataFormatUpsertOne {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.SetSchema(v)
	})
}

// UpdateSchema sets the "schema" field to the value that was provided on create.
func (u *MetadataFormatUpsertOne) UpdateSchema() *MetadataFormatUpsertOne {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.UpdateSchema()
	})
}

// SetNamespace sets the "namespace" field.
func (u *MetadataFormatUpsertOne) SetNamespace(v string) *MetadataFormatUpsertOne {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.SetNamespace(v)
	})
}

// UpdateNamespace sets the "namespace" field to the value that was provided on create.
func (u *MetadataFormatUpsertOne) UpdateNamespace() *MetadataFormatUpsertOne {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.UpdateNamespace()
	})
}

// Exec executes the query.
func (u *MetadataFormatUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MetadataFormatCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MetadataFormatUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *MetadataFormatUpsertOne) ID(ctx context.Context) (id string, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: MetadataFormatUpsertOne.ID is not supported by MySQL driver. Use MetadataFormatUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *MetadataFormatUpsertOne) IDX(ctx context.Context) string {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// MetadataFormatCreateBulk is the builder for creating many MetadataFormat entities in bulk.
type MetadataFormatCreateBulk struct {
	config
	builders []*MetadataFormatCreate
	conflict []sql.ConflictOption
}

// Save creates the MetadataFormat entities in the database.
func (mfcb *MetadataFormatCreateBulk) Save(ctx context.Context) ([]*MetadataFormat, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mfcb.builders))
	nodes := make([]*MetadataFormat, len(mfcb.builders))
	mutators := make([]Mutator, len(mfcb.builders))
	for i := range mfcb.builders {
		func(i int, root context.Context) {
			builder := mfcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MetadataFormatMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mfcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mfcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mfcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mfcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mfcb *MetadataFormatCreateBulk) SaveX(ctx context.Context) []*MetadataFormat {
	v, err := mfcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mfcb *MetadataFormatCreateBulk) Exec(ctx context.Context) error {
	_, err := mfcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mfcb *MetadataFormatCreateBulk) ExecX(ctx context.Context) {
	if err := mfcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.MetadataFormat.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MetadataFormatUpsert) {
//			SetSchema(v+v).
//		}).
//		Exec(ctx)
func (mfcb *MetadataFormatCreateBulk) OnConflict(opts ...sql.ConflictOption) *MetadataFormatUpsertBulk {
	mfcb.conflict = opts
	return &MetadataFormatUpsertBulk{
		create: mfcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.MetadataFormat.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mfcb *MetadataFormatCreateBulk) OnConflictColumns(columns ...string) *MetadataFormatUpsertBulk {
	mfcb.conflict = append(mfcb.conflict, sql.ConflictColumns(columns...))
	return &MetadataFormatUpsertBulk{
		create: mfcb,
	}
}

// MetadataFormatUpsertBulk is the builder for "upsert"-ing
// a bulk of MetadataFormat nodes.
type MetadataFormatUpsertBulk struct {
	create *MetadataFormatCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.MetadataFormat.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(metadataformat.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MetadataFormatUpsertBulk) UpdateNewValues() *MetadataFormatUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(metadataformat.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.MetadataFormat.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *MetadataFormatUpsertBulk) Ignore() *MetadataFormatUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MetadataFormatUpsertBulk) DoNothing() *MetadataFormatUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MetadataFormatCreateBulk.OnConflict
// documentation for more info.
func (u *MetadataFormatUpsertBulk) Update(set func(*MetadataFormatUpsert)) *MetadataFormatUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MetadataFormatUpsert{UpdateSet: update})
	}))
	return u
}

// SetSchema sets the "schema" field.
func (u *MetadataFormatUpsertBulk) SetSchema(v string) *MetadataFormatUpsertBulk {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.SetSchema(v)
	})
}

// UpdateSchema sets the "schema" field to the value that was provided on create.
func (u *MetadataFormatUpsertBulk) UpdateSchema() *MetadataFormatUpsertBulk {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.UpdateSchema()
	})
}

// SetNamespace sets the "namespace" field.
func (u *MetadataFormatUpsertBulk) SetNamespace(v string) *MetadataFormatUpsertBulk {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.SetNamespace(v)
	})
}

// UpdateNamespace sets the "namespace" field to the value that was provided on create.
func (u *MetadataFormatUpsertBulk) UpdateNamespace() *MetadataFormatUpsertBulk {
	return u.Update(func(s *MetadataFormatUpsert) {
		s.UpdateNamespace()
	})
}

// Exec executes the query.
func (u *MetadataFormatUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the MetadataFormatCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MetadataFormatCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MetadataFormatUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
