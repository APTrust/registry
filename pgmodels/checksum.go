package pgmodels

import (
	"time"
	//"github.com/APTrust/registry/common"
	//"github.com/APTrust/registry/constants"
)

type Checksum struct {
	ID            int64     `json:"id" form:"id" pg:"id"`
	Algorithm     string    `json:"algorithm" form:"algorithm" pg:"algorithm"`
	DateTime      time.Time `json:"datetime" form:"datetime" pg:"datetime"`
	Digest        string    `json:"digest" form:"digest" pg:"digest"`
	GenericFileID int64     `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	CreatedAt     time.Time `json:"created_at" form:"created_at" pg:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" form:"updated_at" pg:"updated_at"`

	GenericFile *GenericFile `json:"-" pg:"rel:has-one"`
}
