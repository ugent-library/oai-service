package oai

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go/micro"
	"github.com/ugent-library/oai-service/repositories"
)

var Config = micro.Config{
	Name:        "OaiService",
	Version:     "0.0.1",
	Description: "OAI Provider management",
}

type Service struct {
	repo *repositories.Repo
}

func NewService(repo *repositories.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type AddMetadataFormatRequest struct {
	Prefix    string `json:"prefix"`
	Schema    string `json:"schema"`
	Namespace string `json:"namespace"`
}

type AddMetadataFormatResponse struct {
}

func (s *Service) AddMetadataFormat(r micro.Request) {
	req := AddMetadataFormatRequest{}
	if err := json.Unmarshal(r.Data(), &req); err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	err := s.repo.AddMetadataFormat(context.Background(),
		req.Prefix,
		req.Schema,
		req.Namespace,
	)
	if err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	r.RespondJSON(&AddMetadataFormatResponse{})
}

type AddSetRequest struct {
	Spec        string `json:"spec"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddSetResponse struct {
}

func (s *Service) AddSet(r micro.Request) {
	req := AddSetRequest{}
	if err := json.Unmarshal(r.Data(), &req); err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	err := s.repo.AddSet(context.Background(),
		req.Spec,
		req.Name,
		req.Description,
	)
	if err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	r.RespondJSON(&AddSetResponse{})
}

type AddRecordRequest struct {
	Identifier     string   `json:"identifier"`
	MetadataPrefix string   `json:"metadata_prefix"`
	Metadata       string   `json:"metadata"`
	SetSpecs       []string `json:"set_specs"`
}

type AddRecordResponse struct {
}

func (s *Service) AddRecord(r micro.Request) {
	req := AddRecordRequest{}
	if err := json.Unmarshal(r.Data(), &req); err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	err := s.repo.AddRecord(context.Background(),
		req.Identifier,
		req.MetadataPrefix,
		req.Metadata,
		req.SetSpecs,
	)
	if err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	r.RespondJSON(&AddRecordResponse{})
}

type DeleteRecordRequest struct {
	Identifier string `json:"identifier"`
}

type DeleteRecordResponse struct {
}

func (s *Service) DeleteRecord(r micro.Request) {
	req := DeleteRecordRequest{}
	if err := json.Unmarshal(r.Data(), &req); err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	err := s.repo.DeleteRecord(context.Background(),
		req.Identifier,
	)
	if err != nil {
		r.Error("500", err.Error(), nil)
		return
	}
	r.RespondJSON(&DeleteRecordResponse{})
}
