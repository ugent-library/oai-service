// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/ugent-library/oai-service/ent/item"
)

// Item is the model entity for the Item schema.
type Item struct {
	config
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ItemQuery when eager-loading is set.
	Edges        ItemEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ItemEdges holds the relations/edges for other nodes in the graph.
type ItemEdges struct {
	// Records holds the value of the records edge.
	Records []*Record `json:"records,omitempty"`
	// Sets holds the value of the sets edge.
	Sets []*Set `json:"sets,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// RecordsOrErr returns the Records value or an error if the edge
// was not loaded in eager-loading.
func (e ItemEdges) RecordsOrErr() ([]*Record, error) {
	if e.loadedTypes[0] {
		return e.Records, nil
	}
	return nil, &NotLoadedError{edge: "records"}
}

// SetsOrErr returns the Sets value or an error if the edge
// was not loaded in eager-loading.
func (e ItemEdges) SetsOrErr() ([]*Set, error) {
	if e.loadedTypes[1] {
		return e.Sets, nil
	}
	return nil, &NotLoadedError{edge: "sets"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Item) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case item.FieldID:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Item fields.
func (i *Item) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case item.FieldID:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[j])
			} else if value.Valid {
				i.ID = value.String
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Item.
// This includes values selected through modifiers, order, etc.
func (i *Item) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// QueryRecords queries the "records" edge of the Item entity.
func (i *Item) QueryRecords() *RecordQuery {
	return NewItemClient(i.config).QueryRecords(i)
}

// QuerySets queries the "sets" edge of the Item entity.
func (i *Item) QuerySets() *SetQuery {
	return NewItemClient(i.config).QuerySets(i)
}

// Update returns a builder for updating this Item.
// Note that you need to call Item.Unwrap() before calling this method if this Item
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Item) Update() *ItemUpdateOne {
	return NewItemClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Item entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Item) Unwrap() *Item {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Item is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Item) String() string {
	var builder strings.Builder
	builder.WriteString("Item(")
	builder.WriteString(fmt.Sprintf("id=%v", i.ID))
	builder.WriteByte(')')
	return builder.String()
}

// Items is a parsable slice of Item.
type Items []*Item
