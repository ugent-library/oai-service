// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// AddMetadataFormat implements addMetadataFormat operation.
	//
	// Add a metadata format.
	//
	// POST /add-metadata-format
	AddMetadataFormat(ctx context.Context, req *AddMetadataFormatRequest) error
	// AddRecordMetadata implements addRecordMetadata operation.
	//
	// Add record metadata.
	//
	// POST /add-record-metadata
	AddRecordMetadata(ctx context.Context, req *AddRecordMetadataRequest) error
	// AddRecordSets implements addRecordSets operation.
	//
	// Add record sets.
	//
	// POST /add-record-sets
	AddRecordSets(ctx context.Context, req *AddRecordSetsRequest) error
	// AddSet implements addSet operation.
	//
	// Add a set.
	//
	// POST /add-set
	AddSet(ctx context.Context, req *AddSetRequest) error
	// DeleteRecord implements deleteRecord operation.
	//
	// Delete a record.
	//
	// POST /delete-record
	DeleteRecord(ctx context.Context, req *DeleteRecordRequest) error
	// NewError creates *ErrorStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *ErrorStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
