package web

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type Request struct {
	PathAndQuery string
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
	pathAndQuery := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		pathAndQuery = c.Request.URL.Path + "?" + c.Request.URL.RawQuery
	}
	csrfToken, _ := c.Get(constants.CSRFTokenName)
	req := &Request{
		PathAndQuery: pathAndQuery,
		CurrentUser:  currentUser,
		GinContext:   c,
		Auth:         auth.(*middleware.ResourceAuthorization),
		TemplateData: gin.H{
			"CurrentUser":           currentUser,
			"flash":                 flash,
			constants.CSRFTokenName: csrfToken,
		},
	}
	helpers.DeleteFlashCookie(c)
	return req
}

// GetFilterCollection returns a collection of filters the user
// wants to apply to an index/list request. These come from the
// query string. Call the ToQuery() method of the returned
// FilterCollection to translate query string params to SQL.
func (req *Request) GetFilterCollection() *pgmodels.FilterCollection {
	allowedFilters := pgmodels.FiltersFor(req.Auth.ResourceType)
	fc := pgmodels.NewFilterCollection()
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

func (req *Request) LoadResourceList(items interface{}, orderBy string, ffConstructor forms.FilterFormConstructor) error {
	// Ensure that items is a pointer to a slice of pointers, so we don't
	// get a panic in call to Elem() below.
	if items == nil || !strings.HasPrefix(reflect.TypeOf(items).String(), "*[]*pgmodels.") {
		common.Context().Log.Error().Msgf("Request.LoadResourceList: Param items should be pointer to slice of pointers.")
		return common.ErrInvalidParam
	}

	filterCollection := req.GetFilterCollection()
	query, err := filterCollection.ToQuery()
	if err != nil {
		return err
	}
	if !req.CurrentUser.IsAdmin() {
		query.Where("institution_id", "=", req.CurrentUser.InstitutionID)
		if reflect.ValueOf(items).Elem().Type() == reflect.TypeOf([]*pgmodels.AlertView{}) {
			query.Where("user_id", "=", req.CurrentUser.ID)
		}
	}
	query.OrderBy(orderBy)
	pager, err := NewPager(req.GinContext, req.PathAndQuery, 20)
	if err != nil {
		return err
	}
	query.Offset(pager.QueryOffset).Limit(pager.PerPage)
	err = query.Select(items)
	if err != nil {
		return err
	}
	count, err := query.Count(items)
	if err != nil {
		return err
	}
	pager.SetCounts(count, reflect.ValueOf(items).Elem().Len())

	form, err := ffConstructor(filterCollection, req.CurrentUser)

	req.TemplateData["items"] = items
	req.TemplateData["pager"] = pager
	req.TemplateData["filterForm"] = form

	return err
}
