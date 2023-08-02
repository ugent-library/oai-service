// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AddMetadataFormat implements addMetadataFormat operation.
//
// Add a metadata format.
//
// POST /add-metadata-format
func (UnimplementedHandler) AddMetadataFormat(ctx context.Context, req *AddMetadataFormatRequest) error {
	return ht.ErrNotImplemented
}

// AddRecordMetadata implements addRecordMetadata operation.
//
// Add record metadata.
//
// POST /add-record-metadata
func (UnimplementedHandler) AddRecordMetadata(ctx context.Context, req *AddRecordMetadataRequest) error {
	return ht.ErrNotImplemented
}

// AddRecordSets implements addRecordSets operation.
//
// Add record sets.
//
// POST /add-record-sets
func (UnimplementedHandler) AddRecordSets(ctx context.Context, req *AddRecordSetsRequest) error {
	return ht.ErrNotImplemented
}

// AddSet implements addSet operation.
//
// Add a set.
//
// POST /add-set
func (UnimplementedHandler) AddSet(ctx context.Context, req *AddSetRequest) error {
	return ht.ErrNotImplemented
}

// DeleteRecord implements deleteRecord operation.
//
// Delete a record.
//
// POST /delete-record
func (UnimplementedHandler) DeleteRecord(ctx context.Context, req *DeleteRecordRequest) error {
	return ht.ErrNotImplemented
}

// NewError creates *ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *ErrorStatusCode) {
	r = new(ErrorStatusCode)
	return r
}
