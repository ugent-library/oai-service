package api

import (
	"context"

	"github.com/ugent-library/oai-service/repositories"
)

type Service struct {
	repo *repositories.Repo
}

func NewService(repo *repositories.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddRecord(ctx context.Context, req *AddRecordRequest) (*AddRecordResponse, error) {
	err := s.repo.AddRecord(ctx,
		req.Identifier,
		req.MetadataPrefix,
		req.Metadata,
		req.SetSpecs,
	)
	if err != nil {
		return nil, err
	}
	return &AddRecordResponse{}, nil
}

func (s *Service) NewError(ctx context.Context, err error) *ErrorStatusCode {
	return &ErrorStatusCode{
		StatusCode: 500,
		Response: Error{
			Code:    500,
			Message: err.Error(),
		},
	}
}
