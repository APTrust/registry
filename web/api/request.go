package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/stew/slice"
)

type Request struct {
	PathAndQuery string                            `json:"pathAndQuery"`
	CurrentUser  *pgmodels.User                    `json:"currentUser"`
	GinContext   *gin.Context                      `json:"-"`
	Auth         *middleware.ResourceAuthorization `json:"resourceAuth"`
	Error        error                             `json:"error"`
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

// VailidateFilters returns an error if the query string contains invalid
// or unrecognized values. Even though GetFilterCollection below applies
// only valid filters, the caller should know when we are ignoring their
// filters, so they don't mistakenly accept invalid results.
//
// This problem appeared in integration testing with preservation services,
// when some Registry calls included old Pharos filter params that are no
// longer valid in Registry. Registry would silently ignore those filters
// and return *unfiltered* results, which preserv would then act on.
//
// It's much better to fail and force the developer (that jerk!) to fix
// the issues in preservation services.
func (req *Request) ValidateFilters() error {
	allowedFilters := pgmodels.FiltersFor(req.Auth.ResourceType)
	allowedParams := append(allowedFilters, "sort", "page", "per_page")
	invalid := make([]string, 0)
	for paramName, _ := range req.GinContext.Request.URL.Query() {
		if !slice.Contains(allowedParams, paramName) {
			invalid = append(invalid, paramName)
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("Invalid query params: %s", strings.Join(invalid, ", "))
	}
	return nil
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
	for _, value := range req.GinContext.QueryArray("sort") {
		fc.AddOrderBy(value)
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

// LoadResourceList loads a list of resources for an index page.
// Param items should be a pointer to a slice of the type of item
// you want to load (GenericFile, Institution, etc.). Params
// orderByColumn and direction indicate a default sort order to be
// applied if the request did not explicitly include a sort order.
// (I.e. no sort=column__direction on the query string.)
func (req *Request) LoadResourceList(items interface{}, orderByColumn, direction string) (*common.Pager, error) {
	// Ensure that items is a pointer to a slice of pointers, so we don't
	// get a panic in call to Elem() below.
	if items == nil || !strings.HasPrefix(reflect.TypeOf(items).String(), "*[]*pgmodels.") {
		common.Context().Log.Error().Msgf("Request.LoadResourceList: Param items should be pointer to slice of pointers.")
		return nil, common.ErrInvalidParam
	}

	err := req.ValidateFilters()
	if err != nil {
		return nil, err
	}

	filterCollection := req.GetFilterCollection()
	query, err := filterCollection.ToQuery()
	if err != nil {
		return nil, err
	}
	if !req.CurrentUser.IsAdmin() {
		query.Where("institution_id", "=", req.CurrentUser.InstitutionID)
		objType := reflect.ValueOf(items).Elem().Type()
		if objType == reflect.TypeOf([]*pgmodels.AlertView{}) || objType == reflect.TypeOf([]*pgmodels.Alert{}) {
			query.Where("user_id", "=", req.CurrentUser.ID)
		}
	}
	if !filterCollection.HasExplicitSorting() {
		query.OrderBy(orderByColumn, direction)
	}
	pager, err := common.NewPager(req.GinContext, req.PathAndQuery, 20)
	if err != nil {
		return nil, err
	}
	query.Offset(pager.QueryOffset).Limit(pager.PerPage)

	// This sucks. Maybe there's a way to call the underlying
	// type's select method, because that would handle this.
	if reflect.ValueOf(items).Elem().Type() == reflect.TypeOf([]*pgmodels.GenericFile{}) {
		err = query.Relations("Checksums", "PremisEvents", "StorageRecords").Select(items)
	} else {
		err = query.Select(items)
	}
	if err != nil {
		return nil, err
	}

	var count int
	if pgmodels.CanCountFromView(query, items) {
		common.Context().Log.Info().Msgf("API: Using view to count query '%s'", query.WhereClause())
		count, err = pgmodels.GetCountFromView(query, items)
	} else {
		common.Context().Log.Info().Msgf("API: Using standard count query for '%s'", query.WhereClause())
		count, err = query.Count(items)
	}

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
	if req.Auth.ResourceID != resourceID {
		msg += fmt.Sprintf("URL says resource ID %d, but JSON says %d. ", req.Auth.ResourceID, resourceID)
	}
	if req.Auth.ResourceInstID != instID {
		msg += fmt.Sprintf("URL says institution ID %d, but JSON says %d. ", req.Auth.ResourceInstID, instID)
	}
	if len(msg) > 0 {
		common.Context().Log.Error().Msgf("Illegal update. User %s, %s:  %s", req.CurrentUser.Email, req.GinContext.FullPath(), msg)
		return common.ErrIDMismatch
	}
	return nil
}

// ToJson returns the request object as JSON (minus the gin context
// object). This is primarily for interactive debugging. Param pretty
// is for pretty printing.
func (req *Request) ToJson(pretty bool) (string, error) {
	var data []byte
	var err error
	if pretty == true {
		data, err = json.MarshalIndent(req, "", "  ")
	} else {
		data, err = json.Marshal(req)
	}
	return string(data), err
}
