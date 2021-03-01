package controllers

import (
	"strconv"

	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
)

type Req struct {
	CurrentUser  *models.User
	DataStore    *models.DataStore
	ID           int64
	TemplateData gin.H
}

func NewRequest(c *gin.Context) *Req {
	currentUser := helpers.CurrentUser(c)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	return &Req{
		CurrentUser:  currentUser,
		DataStore:    models.NewDataStore(currentUser),
		ID:           id,
		TemplateData: gin.H{"CurrentUser": currentUser},
	}
}
