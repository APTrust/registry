package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type ResourceAuthorization struct {
	c              *gin.Context
	Handler        string
	ResourceID     int64
	ResourceInstID int64
	ResourceType   string
	Permission     constants.Permission
	Checked        bool
	Approved       bool
	Error          error
}

func AuthorizeResource(c *gin.Context) *ResourceAuthorization {
	r := &ResourceAuthorization{c: c}
	r.run()
	return r
}

func (r *ResourceAuthorization) run() {
	r.getPermissionType()
	r.readRequestIds()
	if r.Error == nil {
		r.checkPermission()
	}
}

func (r *ResourceAuthorization) getPermissionType() {
	nameParts := strings.Split(r.c.HandlerName(), ".")
	if len(nameParts) > 1 {
		r.Handler = nameParts[len(nameParts)-1]
		if authMeta, ok := AuthMap[r.Handler]; ok {
			r.Permission = authMeta.Permission
			r.ResourceType = authMeta.ResourceType
		} else {
			r.Error = common.ErrResourcePermission
		}
	}
}

func (r *ResourceAuthorization) checkPermission() {
	currentUser := r.CurrentUser()
	r.Approved = currentUser != nil && currentUser.HasPermission(r.Permission, r.ResourceInstID)
	r.Checked = true
}

func (r *ResourceAuthorization) readRequestIds() {
	r.ResourceID = r.idFromRequest("id")
	r.ResourceInstID = r.idFromRequest("institution_id")
	if strings.HasPrefix(r.Handler, "Institution") {
		r.ResourceInstID = r.ResourceID
	}

	// TODO: Consider forcing institution_id = User.InstitutionID
	// on requests where user is not admin: New, Create, Index.

	if r.ResourceID == int64(0) {
		r.ResourceInstID, r.Error = pgmodels.InstIDFor(r.ResourceType, r.ResourceID)
	}
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

func (r *ResourceAuthorization) CurrentUser() *pgmodels.User {
	if currentUser, ok := r.c.Get("CurrentUser"); ok && currentUser != nil {
		return currentUser.(*pgmodels.User)
	}
	return nil
}

// GetError returns an error message with detailed information.
// This is primarily for logging.
func (r *ResourceAuthorization) GetError() string {
	user, exists := r.c.Get("CurrentUser")
	email := "<user not signed in>"
	if exists && user != nil {
		email = user.(*pgmodels.User).Email
	}
	return fmt.Sprintf("ResourceAuth: User %s, Remote IP: %s, Handler: %s, ResourceID: %d, InstID: %d, Path: %s, Permission: %s, Error: %s", email, r.c.Request.RemoteAddr, r.c.HandlerName(), r.ResourceID, r.ResourceInstID, r.c.FullPath(), r.Permission, r.Error.Error())
}
