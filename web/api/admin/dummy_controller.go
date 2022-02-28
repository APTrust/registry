//go:build !test
// +build !test

package admin_api

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/web/api"
	"github.com/gin-gonic/gin"
)

// PrepareFileDelete is not implemented, except in test and
// integration builds.
//
// POST /admin-api/v3/prepare_file_delete/:id
func PrepareFileDelete(c *gin.Context) {
	api.AbortIfError(c, common.ErrNotSupported)
	return
}

// PrepareObjectDelete is not implemented, except in test and
// integration builds.
//
// POST /admin-api/v3/prepare_object_delete/:id
func PrepareObjectDelete(c *gin.Context) {
	api.AbortIfError(c, common.ErrNotSupported)
	return
}
