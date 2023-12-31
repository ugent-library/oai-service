package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Set struct {
	ent.Schema
}

func (Set) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("spec").
			Unique(),
		field.String("name"),
		field.String("description").Optional(),
	}
}

func (Set) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("items", Item.Type).
			Ref("sets"),
	}
}

func (Set) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("spec").
			Annotations(
				entsql.OpClass("text_pattern_ops"),
			),
	}
}
