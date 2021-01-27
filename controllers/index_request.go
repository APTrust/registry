package controllers

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type IndexRequest struct {
	Context  *gin.Context
	Items    interface{}
	Query    *models.Query
	Template string
}

func NewIndexRequest(c *gin.Context, items interface{}, query *models.Query, template string) *IndexRequest {
	return &IndexRequest{
		Context:  c,
		Items:    items,
		Query:    query,
		Template: template,
	}
}

func (r *IndexRequest) Process() {
	templateData := gin.H{}
	template := r.Template
	status := http.StatusOK
	err := models.Select(r.Items, r.Query)
	if err != nil && err != pg.ErrNoRows {
		status = http.StatusBadRequest
		templateData["error"] = err.Error()
		template = "errors/show.html"
	} else {
		templateData = helpers.TemplateVars(r.Context)
		templateData["items"] = r.Items
		fmt.Println(r.Items)
	}
	r.Context.HTML(status, template, templateData)
}
