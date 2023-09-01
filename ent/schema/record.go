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
		field.String("metadata_format_id"),
		field.String("item_id"),
		field.String("metadata").
			Optional().
			Nillable().
			Comment("A record with NULL metadata is considered deleted."),
		field.Time("datestamp").
			Default(time.Now),
	}
}

func (Record) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("metadata_format", MetadataFormat.Type).
			Field("metadata_format_id").
			Ref("records").
			Unique().
			Required(),
		edge.From("item", Item.Type).
			Field("item_id").
			Ref("records").
			Unique().
			Required(),
	}
}

func (Record) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("metadata_format_id", "item_id").
			Unique(),
		index.Fields("datestamp"),
	}
}
