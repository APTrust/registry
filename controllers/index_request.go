package controllers

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type IndexRequest struct {
	GinCtx         *gin.Context
	AllowedFilters []string
	FilterInstID   bool
	Status         int
	Template       string
	TemplateData   gin.H
	currentUser    *models.User
}

func NewIndexRequest(c *gin.Context, allowedFilters []string, filterInstID bool, template string) *IndexRequest {
	return &IndexRequest{
		GinCtx:         c,
		AllowedFilters: allowedFilters,
		FilterInstID:   filterInstID,
		Status:         http.StatusOK,
		Template:       template,
		TemplateData:   helpers.TemplateVars(c),
		currentUser:    helpers.CurrentUser(c),
	}
}

func (r *IndexRequest) Respond() {
	r.GinCtx.HTML(r.Status, r.Template, r.TemplateData)
}

func (r *IndexRequest) SetError(err error) {
	// ErrNoRows is acceptable in an index request, e.g.
	// when user filters restuls and there are no matches.
	// For other errors, we need to display an error page.
	if err != nil && err != pg.ErrNoRows {
		r.GinCtx.Error(err)
		r.Status = StatusCodeForError(err)
		r.TemplateData["error"] = err.Error()
		r.Template = "errors/show.html"
	}
}

// TODO: return filters as where clause + params
func (r *IndexRequest) GetFilters() (map[string]string, error) {
	filters := make(map[string]string)
	if r.currentUser == nil {
		return nil, common.ErrNotSignedIn
	}
	if r.FilterInstID && !r.currentUser.IsAdmin() {
		filters["institution_id"] = fmt.Sprintf("%d", r.currentUser.InstitutionID)
	}
	for _, key := range r.AllowedFilters {
		value := r.GinCtx.Query(key)
		if value != "" {
			filters[key] = value
		}
	}
	return filters, nil
}

// Call this after gathering results
func (r *IndexRequest) AssertPermissions(models []models.Model) error {
	if r.currentUser == nil {
		return common.ErrNotSignedIn
	}
	for _, obj := range models {
		err := obj.Authorize(r.currentUser, constants.ActionRead)
		if err != nil {
			return err // will be common.ErrPermissionDenied
		}
	}
	return nil
}

// TODO:
//
// Obj can have both Gin context and App context.
//
// Filter parser should allow for complex queries,
// e.g.
//     created_at__gt=2021-02-09
//     user_id__in=[10,11,12]
//     name__starts_with="thom"
//
// Set results list object as []models.Model
// Build query.
// Execute query and store results in results list.
// Check permissions.
//
//   ... if no error...
//
// Set other template vars (e.g. list of institutions)
// Respond.
//

// Helper method for institutions list. Takes int64 inst id to set selected.
// SetTemplateData(key string, value interface)
