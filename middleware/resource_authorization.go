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

// ResourceAuthorization contains information about the current request
// handler, the resource and action being requested, and whether the
// current user is authorized to do what they're trying to do.
type ResourceAuthorization struct {
	ginCtx             *gin.Context
	Handler            string
	ResourceID         int64
	ResourceIdentifier string
	ResourceInstID     int64
	ResourceType       string
	Permission         constants.Permission
	Checked            bool
	Approved           bool
	Error              error
}

// AuthorizeResource returns a ResourceAuthorization struct
// describing what is being authorized and whether the current
// user is allowed to do what they're trying to do.
func AuthorizeResource(c *gin.Context) *ResourceAuthorization {
	r := &ResourceAuthorization{
		ginCtx: c,
	}
	r.init()
	common.ConsoleDebug(r.String())
	return r
}

func (r *ResourceAuthorization) init() {
	if ExemptFromAuth(r.ginCtx) {
		r.Handler = "ExemptHandler"
		r.ResourceType = "Exempt"
		r.Checked = true
		r.Approved = true
		return
	}
	if r.NonAdminIsRequestingAdminAPI() {
		r.Handler = "AdminAPI"
		r.ResourceType = "Forbidden"
		r.Checked = true
		r.Approved = false
		r.Error = common.ErrWrongAPI
		return
	}
	r.getPermissionType()
	if r.Error == nil {
		r.readRequestIds()
	}
	if r.Error == nil {
		r.checkPermission()
	}
}

// getPermissionType figures out the resource type the user
// is requesting and the action they are trying to perform on
// that resource.
//
// HandlerName should be the name of a function in the web
// namespace. URLs are mapped to handlers in registry.go.
// If you see an anonymous handler name like "func1", that usually
// means the user requested a route not defined in registry.go.
// This happens most often when we get a GET request on a route
// that is defined for POST or PUT.
func (r *ResourceAuthorization) getPermissionType() {
	nameParts := strings.Split(r.ginCtx.HandlerName(), ".")
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
	id := r.idFromRequest("id")
	r.ResourceID = id
	// id may be an identifier in the URL (not in post body or query string)
	if r.ResourceID == 0 && r.ginCtx.Param("id") != "" {
		idStr := r.ginCtx.Param("id")
		// Not sure why we sometimes get this leading slash.
		// It's dirty, and we should fix it.
		if strings.HasPrefix(idStr, "/") {
			idStr = idStr[1:]
		}
		r.ResourceIdentifier = idStr
		r.ResourceID, r.Error = r.idFromIdentifier()
	}
	r.ResourceInstID = r.idFromRequest("institution_id")
	if r.ResourceInstID == 0 {
		r.ResourceInstID = r.idFromRequest("InstitutionID")
	}
	if strings.HasPrefix(r.Handler, "Institution") {
		r.ResourceInstID = r.ResourceID
	}

	if r.ResourceID != 0 {
		r.ResourceInstID, r.Error = pgmodels.InstIDFor(r.ResourceType, r.ResourceID)
	} else {

		// If institution ID is nowhere in request and we can't get it from
		// the resource either, force institution ID to the user's own
		// institution for everyone except Sys Admin. All non-sysadmins
		// are restricted to their own institution.
		currentUser := r.CurrentUser()
		if r.ResourceInstID == 0 && !currentUser.IsAdmin() {
			r.ResourceInstID = currentUser.InstitutionID
		}
	}
}

func (r *ResourceAuthorization) idFromRequest(name string) int64 {
	id := r.ginCtx.Param(name)
	if id == "" {
		id = r.ginCtx.Query(name)
	}
	if id == "" {
		id = r.ginCtx.PostForm(name)
	}
	// In routes with *id glob patterns, the router picks up
	// the leading slash in addition to the id. We need to get
	// rid of that to parse the id as an int. This affects API
	// routes only. The drawback to the speedy julienschmidt/httprouter
	// is that you have to inject slow code like this to handle
	// common cases.
	id = strings.TrimPrefix(id, "/")
	idAsInt, _ := strconv.ParseInt(id, 10, 64)
	return idAsInt
}

// To be compatible with the old Pharos API, we need to allow
// users to look up institutions, objects, files and events
// by identifier. For events, identifier is a UUID string. For
// the others, it's a semantic string identifier.
//
// Institution identifiers are domain names. E.g. ncsu.edu.
// Object identifiers follow the pattern domain.edu/object_name
// File identifiers use pattern domain.edu/object_name/path/to/file/in/bag
func (r *ResourceAuthorization) idFromIdentifier() (int64, error) {
	switch r.ResourceType {
	case "Institution":
		return pgmodels.IdForInstIdentifier(r.ResourceIdentifier)
	case "IntellectualObject":
		return pgmodels.IdForObjIdentifier(r.ResourceIdentifier)
	case "GenericFile":
		return pgmodels.IdForFileIdentifier(r.ResourceIdentifier)
	case "PremisEvent":
		return pgmodels.IdForEventIdentifier(r.ResourceIdentifier)
	}
	common.Context().Log.Error().Msgf("Resource auth middleware cannot look up id for identifier '%s' on resource type '%s'", r.ResourceIdentifier, r.ResourceType)
	return 0, common.ErrInvalidParam
}

func (r *ResourceAuthorization) CurrentUser() *pgmodels.User {
	if currentUser, ok := r.ginCtx.Get("CurrentUser"); ok && currentUser != nil {
		return currentUser.(*pgmodels.User)
	}
	return nil
}

// GetError returns an error message with detailed information.
// This is primarily for logging.
func (r *ResourceAuthorization) GetError() string {
	return fmt.Sprintf("ResourceAuth Error: %s", r.String())
}

// GetNotAuthorizedMessage returns a message describing what was not
// authorized, and for whom.
func (r *ResourceAuthorization) GetNotAuthorizedMessage() string {
	return fmt.Sprintf("Not Authorized: %s", r.String())
}

// NonAdminIsRequestingAdminAPI returns true if a non-admin user is requesting
// a resource from the admin API. Although the admin and member APIs share
// some common handlers, we want to force members to access features through
// member-api endpoints.
//
// This test is a shortcut that allows us to skip more complicated checks.
func (r *ResourceAuthorization) NonAdminIsRequestingAdminAPI() bool {
	currentUser := r.CurrentUser()
	return strings.HasPrefix(r.ginCtx.Request.URL.Path, constants.APIPrefixAdmin) && (currentUser == nil || !currentUser.IsAdmin())
}

// String returns this object in string format, suitable for debugging.
func (r *ResourceAuthorization) String() string {
	user, exists := r.ginCtx.Get("CurrentUser")
	email := "<user not signed in>"
	if exists && user != nil {
		email = user.(*pgmodels.User).Email
	}
	errMsg := ""
	if r.Error != nil {
		errMsg = r.Error.Error()
	}
	return fmt.Sprintf("User %s, Remote IP: %s, Handler: %s, ResourceType: %s, ResourceID: %d, InstID: %d, Gin Path: %s, Request Path: %s, Permission: %s, ResourceIdentifier: %s, Error: %s", email, r.ginCtx.Request.RemoteAddr, r.ginCtx.HandlerName(), r.ResourceType, r.ResourceID, r.ResourceInstID, r.ginCtx.FullPath(), r.ginCtx.Request.URL.Path, r.Permission, r.ResourceIdentifier, errMsg)
}
