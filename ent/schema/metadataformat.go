package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type MetadataFormat struct {
	ent.Schema
}

func (MetadataFormat) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("schema"),
		field.String("namespace"),
	}
}

func (MetadataFormat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("records", Record.Type),
	}
}
