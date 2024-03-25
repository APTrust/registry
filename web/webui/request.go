package webui

import (
	"encoding/json"
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
	"github.com/stretchr/stew/slice"
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
	resourceAuth, _ := c.Get("ResourceAuthorization")
	pathAndQuery := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		pathAndQuery = c.Request.URL.Path + "?" + c.Request.URL.RawQuery
	}
	csrfToken, _ := c.Get(constants.CSRFTokenName)
	auth := resourceAuth.(*middleware.ResourceAuthorization)
	req := &Request{
		PathAndQuery: pathAndQuery,
		CurrentUser:  currentUser,
		GinContext:   c,
		Auth:         auth,
		TemplateData: gin.H{
			"CurrentUser":           currentUser,
			"filterChips":           make([]*pgmodels.ParamFilter, 0),
			"filterChipJson":        "",
			"flash":                 flash,
			"showAsModal":           common.IsTrueString(c.Query("modal")),
			"openSubMenu":           ShowOpenSubMenu(auth),
			"currentUrl":            c.Request.URL,
			"pageTitle":             auth.PageTitle,
			constants.CSRFTokenName: csrfToken,
		},
	}
	helpers.DeleteFlashCookie(c)
	return req
}

func ShowOpenSubMenu(auth *middleware.ResourceAuthorization) bool {
	submenuItems := []string{
		"AlertIndex",
		"BillingReportShow",
		"DeletionRequestIndex",
		"GenericFileIndex",
		"GenericFileShow",
		"InstitutionIndex",
		"InstitutionShow",
		"InternalMetadataIndex",
		"PremisEventIndex",
		"NsqShow",
		"NsqInit",
		"NsqAdmin",
	}
	return slice.Contains(submenuItems, auth.Handler)
}

// GetFilterCollection returns a collection of filters the user
// wants to apply to an index/list request. These come from the
// query string. Call the ToQuery() method of the returned
// FilterCollection to translate query string params to SQL.
func (req *Request) GetFilterCollection() *pgmodels.FilterCollection {
	chips := make([]*pgmodels.ParamFilter, 0)
	allowedFilters := pgmodels.FiltersFor(req.Auth.ResourceType)
	fc := pgmodels.NewFilterCollection()
	for _, key := range allowedFilters {
		queryValues := req.GinContext.QueryArray(key)
		filter, _ := fc.Add(key, queryValues)
		if !common.ListIsEmpty(queryValues) {
			chips = append(chips, filter)
		}
	}
	for _, value := range req.GinContext.QueryArray("sort") {
		fc.AddOrderBy(value)
	}
	req.TemplateData["filterChips"] = chips
	chipJson, _ := json.Marshal(chips)
	req.TemplateData["filterChipJson"] = string(chipJson)
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
func (req *Request) LoadResourceList(items interface{}, orderByColumn, direction string, ffConstructor forms.FilterFormConstructor) error {
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
		return err
	}
	query.Offset(pager.QueryOffset).Limit(pager.PerPage)
	err = query.Select(items)
	if err != nil {
		common.Context().Log.Error().Msgf("Error running main query in WebUI LoadResourceItemList. Where = %s. Error = %v", query.WhereClause(), err)
		return err
	}
	var count int
	if pgmodels.CanCountFromView(query, items) {
		common.Context().Log.Info().Msgf("WebUI: Using view to count query '%s'", query.WhereClause())
		count, err = pgmodels.GetCountFromView(query, items)
		if err != nil {
			common.Context().Log.Error().Msgf("Error running count query on view with where clause %s: %v", query.WhereClause(), err)
		}
	} else {
		common.Context().Log.Info().Msgf("WebUI: Using standard count query for '%s'", query.WhereClause())
		count, err = query.Count(items)
		if err != nil {
			common.Context().Log.Error().Msgf("Error running standard count with where clause %s: %v", query.WhereClause(), err)
		}
	}

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
