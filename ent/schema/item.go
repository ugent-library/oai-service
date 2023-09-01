package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Item struct {
	ent.Schema
}

func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("identifier").
			Unique(),
	}
}

func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("records", Record.Type),
		edge.To("sets", Set.Type),
	}
}
