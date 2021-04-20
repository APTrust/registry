package web

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type Request struct {
	CurrentUser  *pgmodels.User
	GinContext   *gin.Context
	Auth         *middleware.ResourceAuthorization
	TemplateData gin.H
	Error        error
}

func NewRequest(c *gin.Context) *Request {
	currentUser := helpers.CurrentUser(c)
	auth, _ := c.Get("ResourceAuthorization")
	return &Request{
		CurrentUser:  currentUser,
		GinContext:   c,
		Auth:         auth.(*middleware.ResourceAuthorization),
		TemplateData: gin.H{"CurrentUser": currentUser},
	}
}

func (req *Request) LoadResource() {
	switch req.Auth.ResourceType {
	case "Institution":
		req.TemplateData["item"], req.Error = pgmodels.InstitutionByID(req.Auth.ResourceID)
	case "User":
		req.TemplateData["item"], req.Error = pgmodels.UserByID(req.Auth.ResourceID)
	case "WorkItem":
		req.TemplateData["item"], req.Error = pgmodels.WorkItemByID(req.Auth.ResourceID)
	default:
		req.Error = common.ErrNotSupported
	}
}

func (req *Request) LoadResourceList() {

}

func (req *Request) InitForm() {

}
