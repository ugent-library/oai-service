// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/ugent-library/oai-service/ent/set"
)

// Set is the model entity for the Set schema.
type Set struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SetQuery when eager-loading is set.
	Edges        SetEdges `json:"edges"`
	selectValues sql.SelectValues
}

// SetEdges holds the relations/edges for other nodes in the graph.
type SetEdges struct {
	// Items holds the value of the items edge.
	Items []*Item `json:"items,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ItemsOrErr returns the Items value or an error if the edge
// was not loaded in eager-loading.
func (e SetEdges) ItemsOrErr() ([]*Item, error) {
	if e.loadedTypes[0] {
		return e.Items, nil
	}
	return nil, &NotLoadedError{edge: "items"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Set) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case set.FieldID, set.FieldName, set.FieldDescription:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Set fields.
func (s *Set) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case set.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				s.ID = value.String
			}
		case set.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case set.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				s.Description = value.String
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Set.
// This includes values selected through modifiers, order, etc.
func (s *Set) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryItems queries the "items" edge of the Set entity.
func (s *Set) QueryItems() *ItemQuery {
	return NewSetClient(s.config).QueryItems(s)
}

// Update returns a builder for updating this Set.
// Note that you need to call Set.Unwrap() before calling this method if this Set
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Set) Update() *SetUpdateOne {
	return NewSetClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Set entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Set) Unwrap() *Set {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Set is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Set) String() string {
	var builder strings.Builder
	builder.WriteString("Set(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("name=")
	builder.WriteString(s.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(s.Description)
	builder.WriteByte(')')
	return builder.String()
}

// Sets is a parsable slice of Set.
type Sets []*Set
