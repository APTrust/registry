package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
)

type ResourceAuthorization struct {
	c              *gin.Context
	Handler        string
	ResourceID     int64
	ResourceInstID int64
	Permission     constants.Permission
	Checked        bool
	Approved       bool
	Error          error
}

func NewResourceAuthorization(c *gin.Context) *ResourceAuthorization {
	r := &ResourceAuthorization{c: c}
	r.run()
	return r
}

func (r *ResourceAuthorization) run() {
	r.getPermissionType()
	r.readRequestIds()
	if r.Error == nil {
		r.getInstitutionID()
	}
	if r.Error == nil {
		r.checkPermission()
	}
}

func (r *ResourceAuthorization) getPermissionType() {
	nameParts := strings.Split(r.c.HandlerName(), ".")
	if len(nameParts) > 1 {
		r.Handler = nameParts[len(nameParts)-1]
		r.Permission = constants.PermissionForHandler[r.Handler]
	}
	if r.Permission == "" {
		r.Error = common.ErrResourcePermission
	}
}

func (r *ResourceAuthorization) getInstitutionID() {
	// Ask DB for institution id of resource.
	// Set ResourceInstID, Error
}

func (r *ResourceAuthorization) checkPermission() {
	// Use User.HasPermission()
	// Set Checked, Approved, Error
}

func (r *ResourceAuthorization) readRequestIds() {
	r.ResourceID = r.idFromRequest("id")
	r.ResourceInstID = r.idFromRequest("institution_id")
}

func (r *ResourceAuthorization) idFromRequest(name string) int64 {
	id := r.c.Param(name)
	if id == "" {
		id = r.c.Query(name)
	}
	if id == "" {
		id = r.c.PostForm(name)
	}
	idAsInt, _ := strconv.ParseInt(id, 10, 64)
	return idAsInt
}

// GetError returns an error message with detailed information.
// This is primarily for logging.
func (r *ResourceAuthorization) GetError() string {
	user, exists := r.c.Get("CurrentUser")
	email := "<user not signed in>"
	if exists && user != nil {
		email = user.(*models.User).Email
	}
	return fmt.Sprintf("ResourceAuth: User %s, Remote IP: %s, Handler: %s, ResourceID: %d, InstID: %d, Path: %s, Permission: %s, Error: %s", email, r.c.Request.RemoteAddr, r.c.HandlerName(), r.ResourceID, r.ResourceInstID, r.c.FullPath(), r.Permission, r.Error.Error())
}
