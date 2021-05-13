package pgmodels

import (
	"time"
)

type Checksum struct {
	ID            int64     `json:"id"`
	Algorithm     string    `json:"algorithm"`
	DateTime      time.Time `json:"datetime" pg:"datetime"`
	Digest        string    `json:"digest"`
	GenericFileID int64     `json:"generic_file_id" pg:"generic_file_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	GenericFile *GenericFile `json:"-" pg:"rel:has-one"`
}
