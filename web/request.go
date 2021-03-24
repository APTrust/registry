package web

import (
	"strconv"

	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type Request struct {
	CurrentUser  *pgmodels.User
	GinContext   *gin.Context
	ResourceID   int64
	TemplateData gin.H
}

func NewRequest(c *gin.Context) *Request {
	currentUser := helpers.CurrentUser(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	return &Request{
		CurrentUser:  currentUser,
		GinContext:   c,
		ResourceID:   id,
		TemplateData: gin.H{"CurrentUser": currentUser},
	}
}

// TODO: Load resource or resource list here.
