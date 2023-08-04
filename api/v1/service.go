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

func (s *Service) AddMetadataFormat(ctx context.Context, req *AddMetadataFormatRequest) error {
	err := s.repo.AddMetadataFormat(ctx,
		req.MetadataPrefix,
		req.Schema,
		req.MetadataNamespace,
	)
	return err
}

func (s *Service) AddSet(ctx context.Context, req *AddSetRequest) error {
	err := s.repo.AddSet(ctx,
		req.SetSpec,
		req.SetName,
		req.SetDescription.Or(""),
	)
	return err
}

func (s *Service) AddRecordSets(ctx context.Context, req *AddRecordSetsRequest) error {
	err := s.repo.AddRecordSets(ctx,
		req.Identifier,
		req.SetSpecs,
	)
	return err
}

func (s *Service) AddRecordMetadata(ctx context.Context, req *AddRecordMetadataRequest) error {
	err := s.repo.AddRecordMetadata(ctx,
		req.Identifier,
		req.MetadataPrefix,
		req.XML,
	)
	return err
}

func (s *Service) DeleteRecord(ctx context.Context, req *DeleteRecordRequest) error {
	err := s.repo.DeleteRecord(ctx,
		req.Identifier,
	)
	return err
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
