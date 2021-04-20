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
	req := &Request{
		CurrentUser:  currentUser,
		GinContext:   c,
		Auth:         auth.(*middleware.ResourceAuthorization),
		TemplateData: gin.H{"CurrentUser": currentUser},
	}
	if req.Auth.ShouldLoadResource() {
		req.LoadResource()
	} else if req.Auth.ShouldLoadList() {
		req.LoadResourceList()
	} else if req.Auth.ShouldLoadNewItem() {
		req.LoadNewItem()
	}
	if req.Error == nil && req.Auth.ShouldBind() {
		c.ShouldBind(req.TemplateData["item"])
	}
	return req
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

func (req *Request) LoadNewItem() {
	switch req.Auth.ResourceType {
	case "Institution":
		req.TemplateData["item"] = pgmodels.Institution{}
	case "User":
		req.TemplateData["item"] = pgmodels.User{}
	case "WorkItem":
		req.TemplateData["item"] = pgmodels.WorkItem{}
	default:
		req.Error = common.ErrNotSupported
	}
}

func (req *Request) LoadResourceList() {
	var query *pgmodels.Query
	var err error
	switch req.Auth.ResourceType {
	case "Institution":
		query, err = req.GetIndexQuery(pgmodels.InstitutionFilters)
		if err != nil {
			req.Error = err
		} else {
			req.TemplateData["items"], req.Error = pgmodels.InstitutionViewSelect(query)
		}
	case "User":
		query, err = req.GetIndexQuery(pgmodels.UserFilters)
		if err != nil {
			req.Error = err
		} else {
			req.TemplateData["items"], req.Error = pgmodels.UserViewSelect(query)
		}
	case "WorkItem":
		query, err = req.GetIndexQuery(pgmodels.WorkItemFilters)
		if err != nil {
			req.Error = err
		} else {
			req.TemplateData["items"], req.Error = pgmodels.WorkItemViewSelect(query)
		}
	default:
		req.Error = common.ErrNotSupported
	}
}

func (req *Request) GetIndexQuery(allowedFilters []string) (*pgmodels.Query, error) {
	fc := NewFilterCollection()
	for _, key := range allowedFilters {
		fc.Add(key, req.GinContext.QueryArray(key))
	}
	return fc.ToQuery()
}
