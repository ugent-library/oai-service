package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Metadata struct {
	ent.Schema
}

func (Metadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "metadata"},
	}
}

func (Metadata) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Int64("record_id"),
		field.Int64("metadata_format_id"),
		field.String("metadata"),
		field.Time("datestamp").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Metadata) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("record", Record.Type).
			Ref("metadata").
			Unique().
			Required().
			Field("record_id"),
		edge.From("metadata_format", MetadataFormat.Type).
			Ref("metadata").
			Unique().
			Required().
			Field("metadata_format_id"),
	}
}

func (Metadata) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("record_id", "metadata_format_id").
			Unique(),
		index.Fields("datestamp"),
	}
}
