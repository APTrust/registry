package webui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /ui_components
func ComponentsIndex(c *gin.Context) {
	templateData := gin.H{
		"suppressSideNav": true,
		"suppressTopNav":  true,
	}
	c.HTML(http.StatusOK, "ui_components/index.html", templateData)
}
