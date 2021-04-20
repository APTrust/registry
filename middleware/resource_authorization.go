package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/stew/slice"
)

// The following request types should load a single resource
// with a specific ID.
var shouldLoadResorce = []string{
	constants.RequestTypeShow,
	constants.RequestTypeEdit,
	constants.RequestTypeUpdate,
	constants.RequestTypeCreate,
	constants.RequestTypeDelete,
	constants.RequestTypeUndelete,
}

// The following request types should try to bind submitted
// data to the resource being created/updated.
var shouldBind = []string{
	constants.RequestTypeEdit,
	constants.RequestTypeUpdate,
	constants.RequestTypeCreate,
	constants.RequestTypeDelete,
	constants.RequestTypeUndelete,
}

// ResourceAuthorization contains information about the current request
// handler, the resource and action being requested, and whether the
// current user is authorized to do what they're trying to do.
type ResourceAuthorization struct {
	ginCtx         *gin.Context
	Handler        string
	RequestType    string
	ResourceID     int64
	ResourceInstID int64
	ResourceType   string
	Permission     constants.Permission
	Checked        bool
	Approved       bool
	Error          error
}

// AuthorizeResource returns a ResourceAuthorization struct
// describing what is being authorized and whether the current
// user is allowed to do what they're trying to do.
func AuthorizeResource(c *gin.Context) *ResourceAuthorization {
	r := &ResourceAuthorization{
		ginCtx:      c,
		RequestType: constants.RequestTypeOther,
	}
	r.init()
	fmt.Println(r)
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
	r.getPermissionType()
	r.setRequestType()
	r.readRequestIds()
	if r.Error == nil {
		r.checkPermission()
	}
}

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

func (r *ResourceAuthorization) setRequestType() {
	for _, reqType := range constants.RequestTypes {
		if strings.HasSuffix(r.Handler, reqType) {
			r.RequestType = reqType
			break
		}
	}
}

func (r *ResourceAuthorization) ShouldLoadResource() bool {
	return slice.Contains(shouldLoadResorce, r.RequestType)
}

func (r *ResourceAuthorization) ShouldLoadList() bool {
	return r.RequestType == constants.RequestTypeIndex
}

func (r *ResourceAuthorization) ShouldLoadNewItem() bool {
	return r.RequestType == constants.RequestTypeNew
}

func (r *ResourceAuthorization) ShouldBind() bool {
	return slice.Contains(shouldBind, r.RequestType)
}

func (r *ResourceAuthorization) checkPermission() {
	currentUser := r.CurrentUser()
	r.Approved = currentUser != nil && currentUser.HasPermission(r.Permission, r.ResourceInstID)
	r.Checked = true
}

func (r *ResourceAuthorization) readRequestIds() {
	r.ResourceID = r.idFromRequest("id")
	r.ResourceInstID = r.idFromRequest("institution_id")
	if r.ResourceInstID == 0 {
		r.ResourceInstID = r.idFromRequest("InstitutionID")
	}
	if strings.HasPrefix(r.Handler, "Institution") {
		r.ResourceInstID = r.ResourceID
	}

	// TODO: Consider forcing institution_id = User.InstitutionID
	// on requests where user is not admin: New, Create, Index.

	if r.ResourceID != 0 {
		r.ResourceInstID, r.Error = pgmodels.InstIDFor(r.ResourceType, r.ResourceID)
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
	idAsInt, _ := strconv.ParseInt(id, 10, 64)
	return idAsInt
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
	return fmt.Sprintf("User %s, Remote IP: %s, Handler: %s, ResourceType: %s, ResourceID: %d, InstID: %d, Path: %s, Permission: %s, Error: %s", email, r.ginCtx.Request.RemoteAddr, r.ginCtx.HandlerName(), r.ResourceType, r.ResourceID, r.ResourceInstID, r.ginCtx.FullPath(), r.Permission, errMsg)
}
