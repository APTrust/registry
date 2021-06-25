package web

import (
	"fmt"

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
	ctx := common.Context()
	flash, _ := c.Get(ctx.Config.Cookies.FlashCookie)
	currentUser := helpers.CurrentUser(c)
	auth, _ := c.Get("ResourceAuthorization")
	req := &Request{
		CurrentUser: currentUser,
		GinContext:  c,
		Auth:        auth.(*middleware.ResourceAuthorization),
		TemplateData: gin.H{
			"CurrentUser": currentUser,
			"flash":       flash,
		},
	}
	helpers.DeleteCookie(c, ctx.Config.Cookies.FlashCookie)
	return req
}

// GetFilterCollection returns a collection of filters the user
// wants to apply to an index/list request. These come from the
// query string. Call the ToQuery() method of the returned
// FilterCollection to translate query string params to SQL.
func (req *Request) GetFilterCollection() *FilterCollection {
	allowedFilters := pgmodels.FiltersFor(req.Auth.ResourceType)
	fc := NewFilterCollection()
	for _, key := range allowedFilters {
		fc.Add(key, req.GinContext.QueryArray(key))
	}
	return fc
}

// BaseURL returns the base of param _url. The base includes the scheme,
// optional port, and hostname. In other words, the URL stripped of path
// and query.
func (req *Request) BaseURL() string {
	scheme := common.Context().Config.HTTPScheme()
	host := req.GinContext.Request.Host // host or host:port
	return fmt.Sprintf("%s://%s", scheme, host)
}
