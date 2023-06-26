// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// MetadataFormatsColumns holds the columns for the "metadata_formats" table.
	MetadataFormatsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "prefix", Type: field.TypeString, Unique: true},
		{Name: "schema", Type: field.TypeString},
		{Name: "namespace", Type: field.TypeString},
	}
	// MetadataFormatsTable holds the schema information for the "metadata_formats" table.
	MetadataFormatsTable = &schema.Table{
		Name:       "metadata_formats",
		Columns:    MetadataFormatsColumns,
		PrimaryKey: []*schema.Column{MetadataFormatsColumns[0]},
	}
	// RecordsColumns holds the columns for the "records" table.
	RecordsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "identifier", Type: field.TypeString},
		{Name: "metadata", Type: field.TypeString, Nullable: true},
		{Name: "deleted", Type: field.TypeBool, Default: false},
		{Name: "datestamp", Type: field.TypeTime},
		{Name: "metadata_format_id", Type: field.TypeInt64},
	}
	// RecordsTable holds the schema information for the "records" table.
	RecordsTable = &schema.Table{
		Name:       "records",
		Columns:    RecordsColumns,
		PrimaryKey: []*schema.Column{RecordsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "records_metadata_formats_records",
				Columns:    []*schema.Column{RecordsColumns[5]},
				RefColumns: []*schema.Column{MetadataFormatsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "record_identifier_metadata_format_id",
				Unique:  true,
				Columns: []*schema.Column{RecordsColumns[1], RecordsColumns[5]},
			},
		},
	}
	// SetsColumns holds the columns for the "sets" table.
	SetsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "spec", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
	}
	// SetsTable holds the schema information for the "sets" table.
	SetsTable = &schema.Table{
		Name:       "sets",
		Columns:    SetsColumns,
		PrimaryKey: []*schema.Column{SetsColumns[0]},
	}
	// SetRecordsColumns holds the columns for the "set_records" table.
	SetRecordsColumns = []*schema.Column{
		{Name: "set_id", Type: field.TypeInt64},
		{Name: "record_id", Type: field.TypeInt64},
	}
	// SetRecordsTable holds the schema information for the "set_records" table.
	SetRecordsTable = &schema.Table{
		Name:       "set_records",
		Columns:    SetRecordsColumns,
		PrimaryKey: []*schema.Column{SetRecordsColumns[0], SetRecordsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "set_records_set_id",
				Columns:    []*schema.Column{SetRecordsColumns[0]},
				RefColumns: []*schema.Column{SetsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "set_records_record_id",
				Columns:    []*schema.Column{SetRecordsColumns[1]},
				RefColumns: []*schema.Column{RecordsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		MetadataFormatsTable,
		RecordsTable,
		SetsTable,
		SetRecordsTable,
	}
)

func init() {
	RecordsTable.ForeignKeys[0].RefTable = MetadataFormatsTable
	SetRecordsTable.ForeignKeys[0].RefTable = SetsTable
	SetRecordsTable.ForeignKeys[1].RefTable = RecordsTable
}
