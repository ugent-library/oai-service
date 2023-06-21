package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MetadataFormat holds the schema definition for the MetadataFormat entity.
type MetadataFormat struct {
	ent.Schema
}

// Fields of the MetadataFormat.
func (MetadataFormat) Fields() []ent.Field {
	return []ent.Field{
		field.String("prefix").Unique(),
		field.String("schema"),
		field.String("namespace"),
	}
}

// Edges of the MetadataFormat.
func (MetadataFormat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("records", Record.Type),
	}
}
