package web

import (
	//"fmt"
	"net/http"
	"strconv"
	"time"

	//"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type DepositReportParams struct {
	InstitutionID int64
	StorageOption string
	UpdatedBefore time.Time
}

func DepositReportShow(c *gin.Context) {
	req := NewRequest(c)
	template := "reports/deposits.html"
	params := getDepositReportParams(c)
	if !req.CurrentUser.IsAdmin() {
		params.InstitutionID = req.CurrentUser.InstitutionID
	}
	deposits, err := pgmodels.DepositStatsSelect(params.InstitutionID, params.StorageOption, params.UpdatedBefore)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["deposits"] = deposits
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// getDeositReportParams parses params from the query string for our
// deposit report. It ignores parse errors for updatedBefore and
// institutionID because these fields can legitimately be empty.
func getDepositReportParams(c *gin.Context) DepositReportParams {
	updatedBefore, _ := time.Parse("2006-01-02", c.Query("updated_at__lteq"))
	institutionID, _ := strconv.ParseInt(c.Query("institution_id"), 10, 64)
	storageOption := c.Query("storage_option")
	return DepositReportParams{
		InstitutionID: institutionID,
		StorageOption: storageOption,
		UpdatedBefore: updatedBefore,
	}
}
