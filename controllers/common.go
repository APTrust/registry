package controllers

import (
	"strconv"

	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type Req struct {
	CurrentUser  *pgmodels.User
	ID           int64
	TemplateData gin.H
}

func NewRequest(c *gin.Context) *Req {
	currentUser := helpers.CurrentUser(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	return &Req{
		CurrentUser:  currentUser,
		ID:           id,
		TemplateData: gin.H{"CurrentUser": currentUser},
	}
}
