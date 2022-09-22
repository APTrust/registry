package webui

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/gin-gonic/gin"
)

// AbortIfError stops request processing and displays an
// error page if param err is not nil. This returns true
// if it actually did have to abort. When it returns true,
// the caller should return to ensure that no further processing
// of the request occurs. If this returns false, there was no error,
// and the caller can continue processing.
func AbortIfError(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err)
		c.Set("err", err)
		_, showErrorInModal := c.Get("showErrorInModal")
		if showErrorInModal {
			ErrorShowModal(c)
		} else {
			ErrorShow(c)
		}
		c.Abort()
		return true
	}
	return false
}

func ErrorShow(c *gin.Context) {
	currentUser, _ := c.Get("CurrentUser")
	templateData := gin.H{
		"CurrentUser": currentUser,
	}

	status := http.StatusInternalServerError
	err, _ := c.Get("err")
	if err != nil {
		templateData["error"] = err.(error).Error()
		status = StatusCodeForError(err.(error))
	}
	if currentUser == nil || err == common.ErrMustCompleteReset {
		templateData["suppressSideNav"] = true
		templateData["suppressTopNav"] = true
	}
	c.HTML(status, "errors/show.html", templateData)
}

func ErrorShowModal(c *gin.Context) {
	status := http.StatusInternalServerError
	err, _ := c.Get("err")
	templateData := gin.H{}
	if err != nil {
		templateData["error"] = err.(error).Error()
		status = StatusCodeForError(err.(error))
	}
	c.HTML(status, "errors/show_modal.html", templateData)
}

// StatusCodeForError returns the http.StatusCode for the specified
// error. If the error doesn't map to a code, this returns 500 by
// default.
func StatusCodeForError(err error) (status int) {
	switch err {
	case common.ErrInvalidLogin:
		status = http.StatusUnauthorized
	case common.ErrAccountDeactivated:
		status = http.StatusForbidden
	case common.ErrPermissionDenied, common.ErrMustCompleteReset:
		status = http.StatusForbidden
	case common.ErrParentRecordNotFound:
		status = http.StatusNotFound
	case common.ErrWrongDataType:
		status = http.StatusBadRequest
	case common.ErrDecodeCookie:
		status = http.StatusBadRequest
	case common.ErrNotSupported:
		status = http.StatusMethodNotAllowed
	case common.ErrInternal:
		status = http.StatusInternalServerError
	case common.ErrPendingWorkItems:
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}
	return status
}
