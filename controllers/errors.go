package controllers

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/helpers"
	"github.com/gin-gonic/gin"
)

func ErrorShow(c *gin.Context) {
	c.HTML(200, "errors/show.html", helpers.TemplateVars(c))
}

// StatusCodeForError returns the http.StatusCode for the specified
// error. If the error doesn't map to a code, this returns 500 by
// default.
func StatusCodeForError(err error) (status int) {
	switch err {
	case common.ErrInvalidLogin:
		status = http.StatusUnauthorized
	case common.ErrAccountDeactivated:
	case common.ErrPermissionDenied:
		status = http.StatusForbidden
	case common.ErrParentRecordNotFound:
		status = http.StatusNotFound
	case common.ErrWrongDataType:
	case common.ErrDecodeCookie:
		status = http.StatusBadRequest
	case common.ErrNotSupported:
		status = http.StatusMethodNotAllowed
	default:
		status = http.StatusInternalServerError
	}
	return status
}
