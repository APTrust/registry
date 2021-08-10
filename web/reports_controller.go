package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

type DepositReportParams struct {
	InstitutionID int64
	StorageOption string
	UpdatedBefore time.Time
}

// DepositReportShow shows the deposits report.
//
// Note that this does not follow the usual pattern for list/show
// pages, where most of the work is done by Request or
// Request.LoadResourceList because this is a reporting query that
// is not running against a basic table or view. The query in
// pgmodels.DepositStats is more complex, so we do a little more manual
// work here.
//
// GET /reports/deposits
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
	filterForm, err := forms.NewDepositReportFilterForm(req.GetFilterCollection(), req.CurrentUser)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["deposits"] = deposits
	req.TemplateData["filterForm"] = filterForm
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
