package middleware

import (
	"github.com/APTrust/registry/constants"
	"github.com/gin-gonic/gin"
)

type ResourceAuthorization struct {
	c              *gin.Context
	ResourceType   string
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
	r.parseURL()
	if r.Error == nil {
		r.getInstitutionID()
	}
	if r.Error == nil {
		r.checkPermission()
	}
}

func (r *ResourceAuthorization) parseURL() {
	// Get resource type from route.
	// Get resource id from route (if available).
	// Set ResourceType, ResourceID, Error
}

func (r *ResourceAuthorization) getInstitutionID() {
	// Ask DB for institution id of resource.
	// Set ResourceInstID, Error
}

func (r *ResourceAuthorization) checkPermission() {
	// Use User.HasPermission()
	// Set Checked, Approved, Error
}
