package models

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type MetadataFormat struct {
	Prefix    string `json:"prefix"`
	Namespace string `json:"namespace"`
	Schema    string `json:"schema"`
}

type Record struct {
	Identifier string    `json:"identifier"`
	Datestamp  time.Time `json:"datestamp"`
	Deleted    bool      `json:"deleted"`
	Metadata   string    `json:"metadata,omitempty"`
	SetSpecs   []string  `json:"setSpecs,omitempty"`
}
