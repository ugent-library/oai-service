package api

import (
	"context"

	"github.com/bufbuild/connect-go"
	oaiv1 "github.com/ugent-library/oai-service/gen/oai/v1"
	"github.com/ugent-library/oai-service/repositories"
)

type Server struct {
	repo *repositories.Repo
}

func NewServer(repo *repositories.Repo) *Server {
	return &Server{repo: repo}
}

func (s *Server) AddMetadataFormat(
	ctx context.Context,
	req *connect.Request[oaiv1.AddMetadataFormatRequest],
) (*connect.Response[oaiv1.AddMetadataFormatResponse], error) {
	err := s.repo.AddMetadataFormat(ctx,
		req.Msg.Prefix,
		req.Msg.Schema,
		req.Msg.Namespace,
	)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&oaiv1.AddMetadataFormatResponse{})
	return res, nil
}

func (s *Server) AddSet(
	ctx context.Context,
	req *connect.Request[oaiv1.AddSetRequest],
) (*connect.Response[oaiv1.AddSetResponse], error) {
	err := s.repo.AddSet(ctx,
		req.Msg.Spec,
		req.Msg.Name,
		req.Msg.Description,
	)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&oaiv1.AddSetResponse{})
	return res, nil
}

func (s *Server) AddRecord(
	ctx context.Context,
	req *connect.Request[oaiv1.AddRecordRequest],
) (*connect.Response[oaiv1.AddRecordResponse], error) {
	err := s.repo.AddRecord(ctx,
		req.Msg.Identifier,
		req.Msg.MetadataPrefix,
		req.Msg.Metadata,
		req.Msg.SetSpecs,
	)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&oaiv1.AddRecordResponse{})
	return res, nil
}

func (s *Server) DeleteRecord(
	ctx context.Context,
	req *connect.Request[oaiv1.DeleteRecordRequest],
) (*connect.Response[oaiv1.DeleteRecordResponse], error) {
	if err := s.repo.DeleteRecord(ctx, req.Msg.Identifier); err != nil {
		return nil, err
	}

	res := connect.NewResponse(&oaiv1.DeleteRecordResponse{})
	return res, nil
}
