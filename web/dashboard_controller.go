package web

import (
	"github.com/gin-gonic/gin"
)

func DashboardShow(c *gin.Context) {
	r := NewRequest(c)
	c.HTML(200, "dashboard/show.html", r.TemplateData)
}
