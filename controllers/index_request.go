package controllers

import (
	"net/http"

	"github.com/APTrust/registry/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type IndexRequest struct {
	Context      *gin.Context
	Status       int
	Template     string
	TemplateData gin.H
}

func NewIndexRequest(c *gin.Context, template string) *IndexRequest {
	return &IndexRequest{
		Context:      c,
		Status:       http.StatusOK,
		Template:     template,
		TemplateData: helpers.TemplateVars(c),
	}
}

func (r *IndexRequest) Respond() {
	r.Context.HTML(r.Status, r.Template, r.TemplateData)
}

func (r *IndexRequest) SetError(err error) {
	// ErrNoRows is acceptable in an index request, e.g.
	// when user filters restuls and there are no matches.
	// For other errors, we need to display an error page.
	if err != nil && err != pg.ErrNoRows {
		r.Status = StatusCodeForError(err)
		r.TemplateData["error"] = err.Error()
		r.Template = "errors/show.html"
	}
}
