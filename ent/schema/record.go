package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Record holds the schema definition for the Record entity.
type Record struct {
	ent.Schema
}

// Fields of the Record.
func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.Int("metadata_format_id"),
		field.String("identifier"),
		field.String("metadata"),
		field.Bool("deleted").
			Default(false),
		field.Time("datestamp").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Record.
func (Record) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metadata_format", MetadataFormat.Type).
			Ref("records").
			Unique().
			Required().
			Field("metadata_format_id"),
		edge.From("sets", Set.Type).
			Ref("records"),
	}
}
