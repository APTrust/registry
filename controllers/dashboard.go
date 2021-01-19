package controllers

import (
	"github.com/APTrust/registry/helpers"
	"github.com/gin-gonic/gin"
)

func DashboardShow(c *gin.Context) {
	c.HTML(200, "dashboard/show.html", helpers.TemplateVars(c))
}
