package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Record struct {
	ent.Schema
}

func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("identifier"),
		field.Bool("deleted").
			Default(false),
	}
}

func (Record) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("metadata", Metadata.Type),
		edge.To("sets", Set.Type),
	}
}

func (Record) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("identifier").
			Unique(),
	}
}
