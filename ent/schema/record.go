package schema

import (
	"time"

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
		field.Int64("metadata_format_id"),
		field.String("identifier"),
		field.String("metadata").
			Optional(),
		field.Bool("deleted").
			Default(false),
		field.Time("datestamp").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

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

func (Record) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("identifier", "metadata_format_id").
			Unique(),
	}
}
