// Code generated by ent, DO NOT EDIT.

package metadataformat

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the metadataformat type in the database.
	Label = "metadata_format"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldMetadataPrefix holds the string denoting the metadata_prefix field in the database.
	FieldMetadataPrefix = "metadata_prefix"
	// FieldSchema holds the string denoting the schema field in the database.
	FieldSchema = "schema"
	// FieldMetadataNamespace holds the string denoting the metadata_namespace field in the database.
	FieldMetadataNamespace = "metadata_namespace"
	// EdgeMetadata holds the string denoting the metadata edge name in mutations.
	EdgeMetadata = "metadata"
	// Table holds the table name of the metadataformat in the database.
	Table = "metadata_formats"
	// MetadataTable is the table that holds the metadata relation/edge.
	MetadataTable = "metadata"
	// MetadataInverseTable is the table name for the Metadata entity.
	// It exists in this package in order to avoid circular dependency with the "metadata" package.
	MetadataInverseTable = "metadata"
	// MetadataColumn is the table column denoting the metadata relation/edge.
	MetadataColumn = "metadata_format_id"
)

// Columns holds all SQL columns for metadataformat fields.
var Columns = []string{
	FieldID,
	FieldMetadataPrefix,
	FieldSchema,
	FieldMetadataNamespace,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the MetadataFormat queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByMetadataPrefix orders the results by the metadata_prefix field.
func ByMetadataPrefix(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMetadataPrefix, opts...).ToFunc()
}

// BySchema orders the results by the schema field.
func BySchema(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSchema, opts...).ToFunc()
}

// ByMetadataNamespace orders the results by the metadata_namespace field.
func ByMetadataNamespace(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMetadataNamespace, opts...).ToFunc()
}

// ByMetadataCount orders the results by metadata count.
func ByMetadataCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMetadataStep(), opts...)
	}
}

// ByMetadata orders the results by metadata terms.
func ByMetadata(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMetadataStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newMetadataStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MetadataInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MetadataTable, MetadataColumn),
	)
}
