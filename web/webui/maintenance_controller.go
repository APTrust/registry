package webui

import (
	"net/http"
	"strconv"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/helpers"
	"github.com/gin-gonic/gin"
)

func MaintenanceIndex(c *gin.Context) {
	status := http.StatusOK
	jsonMessage := "Registry is operating normally."
	jsonError := ""
	if common.Context().Config.MaintenanceMode {
		status = http.StatusServiceUnavailable
		jsonMessage = ""
		jsonError = "APTrust Registry is currently undergoing maintenance."
	}

	format := c.Query("format")
	if format == "json" {
		data := map[string]string{
			"StatusCode": strconv.Itoa(status),
			"Error":      jsonError,
			"Message":    jsonMessage,
		}
		c.JSON(status, data)
		return
	}

	data := gin.H{
		"maintenanceMode": common.Context().Config.MaintenanceMode,
		"suppressSideNav": true,
		"suppressTopNav":  true,
		"timestamp":       helpers.DateTimeUS(time.Now()),
	}
	c.HTML(status, "maintenance/index.html", data)
}
