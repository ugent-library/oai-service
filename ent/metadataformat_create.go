// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

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
}

// SetPrefix sets the "prefix" field.
func (mfc *MetadataFormatCreate) SetPrefix(s string) *MetadataFormatCreate {
	mfc.mutation.SetPrefix(s)
	return mfc
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

// AddRecordIDs adds the "records" edge to the Record entity by IDs.
func (mfc *MetadataFormatCreate) AddRecordIDs(ids ...int) *MetadataFormatCreate {
	mfc.mutation.AddRecordIDs(ids...)
	return mfc
}

// AddRecords adds the "records" edges to the Record entity.
func (mfc *MetadataFormatCreate) AddRecords(r ...*Record) *MetadataFormatCreate {
	ids := make([]int, len(r))
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
	var (
		err  error
		node *MetadataFormat
	)
	if len(mfc.hooks) == 0 {
		if err = mfc.check(); err != nil {
			return nil, err
		}
		node, err = mfc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*MetadataFormatMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = mfc.check(); err != nil {
				return nil, err
			}
			mfc.mutation = mutation
			if node, err = mfc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(mfc.hooks) - 1; i >= 0; i-- {
			if mfc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = mfc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, mfc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*MetadataFormat)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from MetadataFormatMutation", v)
		}
		node = nv
	}
	return node, err
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
	if _, ok := mfc.mutation.Prefix(); !ok {
		return &ValidationError{Name: "prefix", err: errors.New(`ent: missing required field "MetadataFormat.prefix"`)}
	}
	if _, ok := mfc.mutation.Schema(); !ok {
		return &ValidationError{Name: "schema", err: errors.New(`ent: missing required field "MetadataFormat.schema"`)}
	}
	if _, ok := mfc.mutation.Namespace(); !ok {
		return &ValidationError{Name: "namespace", err: errors.New(`ent: missing required field "MetadataFormat.namespace"`)}
	}
	return nil
}

func (mfc *MetadataFormatCreate) sqlSave(ctx context.Context) (*MetadataFormat, error) {
	_node, _spec := mfc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mfc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (mfc *MetadataFormatCreate) createSpec() (*MetadataFormat, *sqlgraph.CreateSpec) {
	var (
		_node = &MetadataFormat{config: mfc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: metadataformat.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: metadataformat.FieldID,
			},
		}
	)
	if value, ok := mfc.mutation.Prefix(); ok {
		_spec.SetField(metadataformat.FieldPrefix, field.TypeString, value)
		_node.Prefix = value
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
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: record.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MetadataFormatCreateBulk is the builder for creating many MetadataFormat entities in bulk.
type MetadataFormatCreateBulk struct {
	config
	builders []*MetadataFormatCreate
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
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mfcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
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
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
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
