package webui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /accessibility_statement
func ShowAccessibilityStatement(c *gin.Context) {
	req := NewRequest(c)
	req.TemplateData["suppressTopNav"] = true
	req.TemplateData["suppressSideNav"] = true
	req.TemplateData["pageTitle"] = "Accessibility Statement"
	c.HTML(http.StatusOK, "notices/accessibility_statement.html", req.TemplateData)
}
