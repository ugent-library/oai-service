// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ugent-library/oai-service/ent/metadata"
	"github.com/ugent-library/oai-service/ent/metadataformat"
	"github.com/ugent-library/oai-service/ent/record"
)

// MetadataCreate is the builder for creating a Metadata entity.
type MetadataCreate struct {
	config
	mutation *MetadataMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetRecordID sets the "record_id" field.
func (mc *MetadataCreate) SetRecordID(i int64) *MetadataCreate {
	mc.mutation.SetRecordID(i)
	return mc
}

// SetMetadataFormatID sets the "metadata_format_id" field.
func (mc *MetadataCreate) SetMetadataFormatID(i int64) *MetadataCreate {
	mc.mutation.SetMetadataFormatID(i)
	return mc
}

// SetContent sets the "content" field.
func (mc *MetadataCreate) SetContent(s string) *MetadataCreate {
	mc.mutation.SetContent(s)
	return mc
}

// SetDatestamp sets the "datestamp" field.
func (mc *MetadataCreate) SetDatestamp(t time.Time) *MetadataCreate {
	mc.mutation.SetDatestamp(t)
	return mc
}

// SetNillableDatestamp sets the "datestamp" field if the given value is not nil.
func (mc *MetadataCreate) SetNillableDatestamp(t *time.Time) *MetadataCreate {
	if t != nil {
		mc.SetDatestamp(*t)
	}
	return mc
}

// SetID sets the "id" field.
func (mc *MetadataCreate) SetID(i int64) *MetadataCreate {
	mc.mutation.SetID(i)
	return mc
}

// SetRecord sets the "record" edge to the Record entity.
func (mc *MetadataCreate) SetRecord(r *Record) *MetadataCreate {
	return mc.SetRecordID(r.ID)
}

// SetMetadataFormat sets the "metadata_format" edge to the MetadataFormat entity.
func (mc *MetadataCreate) SetMetadataFormat(m *MetadataFormat) *MetadataCreate {
	return mc.SetMetadataFormatID(m.ID)
}

// Mutation returns the MetadataMutation object of the builder.
func (mc *MetadataCreate) Mutation() *MetadataMutation {
	return mc.mutation
}

// Save creates the Metadata in the database.
func (mc *MetadataCreate) Save(ctx context.Context) (*Metadata, error) {
	mc.defaults()
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MetadataCreate) SaveX(ctx context.Context) *Metadata {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MetadataCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MetadataCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mc *MetadataCreate) defaults() {
	if _, ok := mc.mutation.Datestamp(); !ok {
		v := metadata.DefaultDatestamp()
		mc.mutation.SetDatestamp(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MetadataCreate) check() error {
	if _, ok := mc.mutation.RecordID(); !ok {
		return &ValidationError{Name: "record_id", err: errors.New(`ent: missing required field "Metadata.record_id"`)}
	}
	if _, ok := mc.mutation.MetadataFormatID(); !ok {
		return &ValidationError{Name: "metadata_format_id", err: errors.New(`ent: missing required field "Metadata.metadata_format_id"`)}
	}
	if _, ok := mc.mutation.Content(); !ok {
		return &ValidationError{Name: "content", err: errors.New(`ent: missing required field "Metadata.content"`)}
	}
	if _, ok := mc.mutation.Datestamp(); !ok {
		return &ValidationError{Name: "datestamp", err: errors.New(`ent: missing required field "Metadata.datestamp"`)}
	}
	if _, ok := mc.mutation.RecordID(); !ok {
		return &ValidationError{Name: "record", err: errors.New(`ent: missing required edge "Metadata.record"`)}
	}
	if _, ok := mc.mutation.MetadataFormatID(); !ok {
		return &ValidationError{Name: "metadata_format", err: errors.New(`ent: missing required edge "Metadata.metadata_format"`)}
	}
	return nil
}

func (mc *MetadataCreate) sqlSave(ctx context.Context) (*Metadata, error) {
	if err := mc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	mc.mutation.id = &_node.ID
	mc.mutation.done = true
	return _node, nil
}

func (mc *MetadataCreate) createSpec() (*Metadata, *sqlgraph.CreateSpec) {
	var (
		_node = &Metadata{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(metadata.Table, sqlgraph.NewFieldSpec(metadata.FieldID, field.TypeInt64))
	)
	_spec.OnConflict = mc.conflict
	if id, ok := mc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := mc.mutation.Content(); ok {
		_spec.SetField(metadata.FieldContent, field.TypeString, value)
		_node.Content = value
	}
	if value, ok := mc.mutation.Datestamp(); ok {
		_spec.SetField(metadata.FieldDatestamp, field.TypeTime, value)
		_node.Datestamp = value
	}
	if nodes := mc.mutation.RecordIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   metadata.RecordTable,
			Columns: []string{metadata.RecordColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(record.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.RecordID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.MetadataFormatIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   metadata.MetadataFormatTable,
			Columns: []string{metadata.MetadataFormatColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(metadataformat.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.MetadataFormatID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Metadata.Create().
//		SetRecordID(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MetadataUpsert) {
//			SetRecordID(v+v).
//		}).
//		Exec(ctx)
func (mc *MetadataCreate) OnConflict(opts ...sql.ConflictOption) *MetadataUpsertOne {
	mc.conflict = opts
	return &MetadataUpsertOne{
		create: mc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Metadata.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mc *MetadataCreate) OnConflictColumns(columns ...string) *MetadataUpsertOne {
	mc.conflict = append(mc.conflict, sql.ConflictColumns(columns...))
	return &MetadataUpsertOne{
		create: mc,
	}
}

type (
	// MetadataUpsertOne is the builder for "upsert"-ing
	//  one Metadata node.
	MetadataUpsertOne struct {
		create *MetadataCreate
	}

	// MetadataUpsert is the "OnConflict" setter.
	MetadataUpsert struct {
		*sql.UpdateSet
	}
)

// SetRecordID sets the "record_id" field.
func (u *MetadataUpsert) SetRecordID(v int64) *MetadataUpsert {
	u.Set(metadata.FieldRecordID, v)
	return u
}

// UpdateRecordID sets the "record_id" field to the value that was provided on create.
func (u *MetadataUpsert) UpdateRecordID() *MetadataUpsert {
	u.SetExcluded(metadata.FieldRecordID)
	return u
}

// SetMetadataFormatID sets the "metadata_format_id" field.
func (u *MetadataUpsert) SetMetadataFormatID(v int64) *MetadataUpsert {
	u.Set(metadata.FieldMetadataFormatID, v)
	return u
}

// UpdateMetadataFormatID sets the "metadata_format_id" field to the value that was provided on create.
func (u *MetadataUpsert) UpdateMetadataFormatID() *MetadataUpsert {
	u.SetExcluded(metadata.FieldMetadataFormatID)
	return u
}

// SetContent sets the "content" field.
func (u *MetadataUpsert) SetContent(v string) *MetadataUpsert {
	u.Set(metadata.FieldContent, v)
	return u
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *MetadataUpsert) UpdateContent() *MetadataUpsert {
	u.SetExcluded(metadata.FieldContent)
	return u
}

// SetDatestamp sets the "datestamp" field.
func (u *MetadataUpsert) SetDatestamp(v time.Time) *MetadataUpsert {
	u.Set(metadata.FieldDatestamp, v)
	return u
}

// UpdateDatestamp sets the "datestamp" field to the value that was provided on create.
func (u *MetadataUpsert) UpdateDatestamp() *MetadataUpsert {
	u.SetExcluded(metadata.FieldDatestamp)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Metadata.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(metadata.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MetadataUpsertOne) UpdateNewValues() *MetadataUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(metadata.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Metadata.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *MetadataUpsertOne) Ignore() *MetadataUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MetadataUpsertOne) DoNothing() *MetadataUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MetadataCreate.OnConflict
// documentation for more info.
func (u *MetadataUpsertOne) Update(set func(*MetadataUpsert)) *MetadataUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MetadataUpsert{UpdateSet: update})
	}))
	return u
}

// SetRecordID sets the "record_id" field.
func (u *MetadataUpsertOne) SetRecordID(v int64) *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.SetRecordID(v)
	})
}

// UpdateRecordID sets the "record_id" field to the value that was provided on create.
func (u *MetadataUpsertOne) UpdateRecordID() *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateRecordID()
	})
}

// SetMetadataFormatID sets the "metadata_format_id" field.
func (u *MetadataUpsertOne) SetMetadataFormatID(v int64) *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.SetMetadataFormatID(v)
	})
}

// UpdateMetadataFormatID sets the "metadata_format_id" field to the value that was provided on create.
func (u *MetadataUpsertOne) UpdateMetadataFormatID() *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateMetadataFormatID()
	})
}

// SetContent sets the "content" field.
func (u *MetadataUpsertOne) SetContent(v string) *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *MetadataUpsertOne) UpdateContent() *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateContent()
	})
}

// SetDatestamp sets the "datestamp" field.
func (u *MetadataUpsertOne) SetDatestamp(v time.Time) *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.SetDatestamp(v)
	})
}

// UpdateDatestamp sets the "datestamp" field to the value that was provided on create.
func (u *MetadataUpsertOne) UpdateDatestamp() *MetadataUpsertOne {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateDatestamp()
	})
}

// Exec executes the query.
func (u *MetadataUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MetadataCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MetadataUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *MetadataUpsertOne) ID(ctx context.Context) (id int64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *MetadataUpsertOne) IDX(ctx context.Context) int64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// MetadataCreateBulk is the builder for creating many Metadata entities in bulk.
type MetadataCreateBulk struct {
	config
	builders []*MetadataCreate
	conflict []sql.ConflictOption
}

// Save creates the Metadata entities in the database.
func (mcb *MetadataCreateBulk) Save(ctx context.Context) ([]*Metadata, error) {
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Metadata, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MetadataMutation)
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
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = mcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MetadataCreateBulk) SaveX(ctx context.Context) []*Metadata {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MetadataCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MetadataCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Metadata.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.MetadataUpsert) {
//			SetRecordID(v+v).
//		}).
//		Exec(ctx)
func (mcb *MetadataCreateBulk) OnConflict(opts ...sql.ConflictOption) *MetadataUpsertBulk {
	mcb.conflict = opts
	return &MetadataUpsertBulk{
		create: mcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Metadata.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (mcb *MetadataCreateBulk) OnConflictColumns(columns ...string) *MetadataUpsertBulk {
	mcb.conflict = append(mcb.conflict, sql.ConflictColumns(columns...))
	return &MetadataUpsertBulk{
		create: mcb,
	}
}

// MetadataUpsertBulk is the builder for "upsert"-ing
// a bulk of Metadata nodes.
type MetadataUpsertBulk struct {
	create *MetadataCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Metadata.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(metadata.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *MetadataUpsertBulk) UpdateNewValues() *MetadataUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(metadata.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Metadata.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *MetadataUpsertBulk) Ignore() *MetadataUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *MetadataUpsertBulk) DoNothing() *MetadataUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the MetadataCreateBulk.OnConflict
// documentation for more info.
func (u *MetadataUpsertBulk) Update(set func(*MetadataUpsert)) *MetadataUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&MetadataUpsert{UpdateSet: update})
	}))
	return u
}

// SetRecordID sets the "record_id" field.
func (u *MetadataUpsertBulk) SetRecordID(v int64) *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.SetRecordID(v)
	})
}

// UpdateRecordID sets the "record_id" field to the value that was provided on create.
func (u *MetadataUpsertBulk) UpdateRecordID() *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateRecordID()
	})
}

// SetMetadataFormatID sets the "metadata_format_id" field.
func (u *MetadataUpsertBulk) SetMetadataFormatID(v int64) *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.SetMetadataFormatID(v)
	})
}

// UpdateMetadataFormatID sets the "metadata_format_id" field to the value that was provided on create.
func (u *MetadataUpsertBulk) UpdateMetadataFormatID() *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateMetadataFormatID()
	})
}

// SetContent sets the "content" field.
func (u *MetadataUpsertBulk) SetContent(v string) *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.SetContent(v)
	})
}

// UpdateContent sets the "content" field to the value that was provided on create.
func (u *MetadataUpsertBulk) UpdateContent() *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateContent()
	})
}

// SetDatestamp sets the "datestamp" field.
func (u *MetadataUpsertBulk) SetDatestamp(v time.Time) *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.SetDatestamp(v)
	})
}

// UpdateDatestamp sets the "datestamp" field to the value that was provided on create.
func (u *MetadataUpsertBulk) UpdateDatestamp() *MetadataUpsertBulk {
	return u.Update(func(s *MetadataUpsert) {
		s.UpdateDatestamp()
	})
}

// Exec executes the query.
func (u *MetadataUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the MetadataCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for MetadataCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *MetadataUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
