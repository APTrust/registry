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
	return req
}

func (req *Request) GetIndexQuery() (*pgmodels.Query, error) {
	allowedFilters, err := req.GetAllowedFilters()
	if err != nil {
		return nil, err
	}
	fc := NewFilterCollection()
	for _, key := range allowedFilters {
		fc.Add(key, req.GinContext.QueryArray(key))
	}
	return fc.ToQuery()
}

func (req *Request) GetAllowedFilters() ([]string, error) {
	switch req.Auth.ResourceType {
	case "Institution":
		return pgmodels.InstitutionFilters, nil
	case "User":
		return pgmodels.UserFilters, nil
	case "WorkItem":
		return pgmodels.WorkItemFilters, nil
	default:
		return nil, common.ErrNotSupported
	}
}
