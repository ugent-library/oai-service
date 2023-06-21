package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Set holds the schema definition for the Set entity.
type Set struct {
	ent.Schema
}

// Fields of the Set.
func (Set) Fields() []ent.Field {
	return []ent.Field{
		field.String("spec").Unique(),
		field.String("name"),
		field.String("description").Optional(),
	}
}

// Edges of the Set.
func (Set) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("records", Record.Type),
	}
}
