package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/APTrust/registry/common"
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
	Error        error
}

func NewRequest(c *gin.Context) *Request {
	currentUser := helpers.CurrentUser(c)
	auth, _ := c.Get("ResourceAuthorization")
	pathAndQuery := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		pathAndQuery = c.Request.URL.Path + "?" + c.Request.URL.RawQuery
	}
	req := &Request{
		PathAndQuery: pathAndQuery,
		CurrentUser:  currentUser,
		GinContext:   c,
		Auth:         auth.(*middleware.ResourceAuthorization),
	}
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

func (req *Request) LoadResourceList(items interface{}, orderBy string) (*common.Pager, error) {
	// Ensure that items is a pointer to a slice of pointers, so we don't
	// get a panic in call to Elem() below.
	if items == nil || !strings.HasPrefix(reflect.TypeOf(items).String(), "*[]*pgmodels.") {
		common.Context().Log.Error().Msgf("Request.LoadResourceList: Param items should be pointer to slice of pointers.")
		return nil, common.ErrInvalidParam
	}

	filterCollection := req.GetFilterCollection()
	query, err := filterCollection.ToQuery()
	if err != nil {
		return nil, err
	}
	if !req.CurrentUser.IsAdmin() {
		query.Where("institution_id", "=", req.CurrentUser.InstitutionID)
		if reflect.ValueOf(items).Elem().Type() == reflect.TypeOf([]*pgmodels.AlertView{}) {
			query.Where("user_id", "=", req.CurrentUser.ID)
		}
	}
	query.OrderBy(orderBy)
	pager, err := common.NewPager(req.GinContext, req.PathAndQuery, 20)
	if err != nil {
		return nil, err
	}
	query.Offset(pager.QueryOffset).Limit(pager.PerPage)
	err = query.Select(items)
	if err != nil {
		return nil, err
	}
	count, err := query.Count(items)
	if err != nil {
		return nil, err
	}
	pager.SetCounts(count, reflect.ValueOf(items).Elem().Len())
	return pager, err
}

// AssertValidIDs returns an error if resource or institution ID in an
// endpoint's URL params don't match the resource/institution ID in the
// JSON of the request body. This is for security. E.g. We don't want
// someone posting to a URL that purports to update one object when
// in fact the JSON will be updating a different object.
func (req *Request) AssertValidIDs(resourceID, instID int64) error {
	msg := ""
	if common.NonZeroAndUnequalInt64(req.Auth.ResourceID, resourceID) {
		msg += fmt.Sprintf("URL says resource ID %d, but JSON says %d. ", req.Auth.ResourceID, resourceID)
	}
	if common.NonZeroAndUnequalInt64(req.Auth.ResourceInstID, instID) {
		msg += fmt.Sprintf("URL says institution ID %d, but JSON says %d. ", req.Auth.ResourceID, resourceID)
	}
	if len(msg) > 0 {
		common.Context().Log.Error().Msgf("Illegal update. User %s, %s:  %s", req.CurrentUser.Email, req.GinContext.FullPath(), msg)
		return common.ErrIDMismatch
	}
	return nil
}
