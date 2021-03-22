package pgmodels

import (
//"github.com/APTrust/registry/common"
//"github.com/APTrust/registry/constants"
)

type StorageRecord struct {
	ID            int64        `json:"id" form:"id" pg:"id"`
	GenericFileID int64        `json:"generic_file_id" form:"generic_file_id" pg:"generic_file_id"`
	URL           string       `json:"url" form:"url" pg:"url"`
	GenericFile   *GenericFile `json:"-" pg:"rel:has-one"`
}
