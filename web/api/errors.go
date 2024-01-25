package api

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type RequestError struct {
	StatusCode int
	Error      string
}

// AbortIfError stops request processing and displays an
// error page if param err is not nil. This returns true
// if it actually did have to abort. When it returns true,
// the caller should return to ensure that no further processing
// of the request occurs. If this returns false, there was no error,
// and the caller can continue processing.
func AbortIfError(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err)
		status := StatusCodeForError(err.(error))
		e := RequestError{
			StatusCode: status,
			Error:      err.Error(),
		}
		common.Context().Log.Err(err)
		c.JSON(status, e)
		c.Abort()
		return true
	}
	return false
}

// StatusCodeForError returns the http.StatusCode for the specified
// error. If the error doesn't map to a code, this returns 500 by
// default.
func StatusCodeForError(err error) (status int) {
	if _, ok := err.(*common.ValidationError); ok {
		return http.StatusBadRequest
	}
	if pgmodels.IsNoRowError(err) {
		return http.StatusNotFound
	}

	switch err {
	case common.ErrInvalidLogin:
		status = http.StatusUnauthorized
	case common.ErrAccountDeactivated:
		status = http.StatusForbidden
	case common.ErrPermissionDenied:
		status = http.StatusForbidden
	case common.ErrParentRecordNotFound:
		status = http.StatusNotFound
	case common.ErrNotSupported:
		status = http.StatusMethodNotAllowed
	case common.ErrInternal:
		status = http.StatusInternalServerError
	case common.ErrPendingWorkItems:
		status = http.StatusConflict
	case common.ErrWrongDataType, common.ErrIDMismatch, common.ErrInstIDChange, common.ErrIdentifierChange,
		common.ErrStorageOptionChange, common.ErrDecodeCookie, common.ErrInvalidObjectID, common.ErrInvalidRequestorID:
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}
	return status
}
