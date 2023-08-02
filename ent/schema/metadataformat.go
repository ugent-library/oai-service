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
		field.Int64("id"),
		field.String("metadata_prefix").Unique(),
		field.String("schema"),
		field.String("metadata_namespace"),
	}
}

func (MetadataFormat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("metadata", Metadata.Type),
	}
}
