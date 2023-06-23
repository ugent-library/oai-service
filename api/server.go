package api

import (
	"context"

	"github.com/bufbuild/connect-go"
	oaiv1 "github.com/ugent-library/oai-service/gen/oai/v1"
	"github.com/ugent-library/oai-service/repository"
)

type Server struct {
	repo *repository.Repo
}

func NewServer(repo *repository.Repo) *Server {
	return &Server{repo: repo}
}

func (s *Server) AddRecord(
	ctx context.Context,
	req *connect.Request[oaiv1.AddRecordRequest],
) (*connect.Response[oaiv1.AddRecordResponse], error) {
	res := connect.NewResponse(&oaiv1.AddRecordResponse{})
	return res, nil
}
