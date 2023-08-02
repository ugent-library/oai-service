package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Set struct {
	ent.Schema
}

func (Set) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("set_spec").Unique(),
		field.String("set_name"),
		field.String("set_description").Optional(),
	}
}

func (Set) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("records", Record.Type).
			Ref("sets"),
	}
}
