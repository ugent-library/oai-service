package models

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Record struct {
	Identifier string    `json:"identifier"`
	Datestamp  time.Time `json:"datestamp"`
	Deleted    bool      `json:"deleted"`
	Metadata   string    `json:"metadata"`
	SetSpecs   []string  `json:"setSpecs"`
}
