package webui

import (
	"net/http"
	"strconv"
	"time"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/stew/slice"
)

type DepositReportParams struct {
	ChartMetric   string
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

	instList := depositInstList(deposits)
	storageOptionsList := depositStorageOptions(deposits)

	req.TemplateData["deposits"] = deposits
	req.TemplateData["isSingleInstitutionReport"] = params.InstitutionID > 0
	req.TemplateData["filterForm"] = filterForm
	req.TemplateData["reportParams"] = params
	req.TemplateData["depositInstitutions"] = instList
	req.TemplateData["depositStorageOptions"] = storageOptionsList
	c.HTML(http.StatusOK, template, req.TemplateData)
}

func depositInstList(deposits []*pgmodels.DepositStats) []string {
	instList := make([]string, 0)
	for _, stats := range deposits {
		if !slice.ContainsString(instList, stats.InstitutionName) {
			instList = append(instList, stats.InstitutionName)
		}
	}
	return instList
}

func depositStorageOptions(deposits []*pgmodels.DepositStats) []string {
	list := make([]string, 0)
	for _, stats := range deposits {
		if !slice.ContainsString(list, stats.StorageOption) {
			list = append(list, stats.StorageOption)
		}
	}
	return list
}

// getDeositReportParams parses params from the query string for our
// deposit report. It ignores parse errors for updatedBefore and
// institutionID because these fields can legitimately be empty.
func getDepositReportParams(c *gin.Context) DepositReportParams {
	updatedBefore, _ := time.Parse("2006-01-02", c.Query("updated_at__lteq"))
	if updatedBefore.IsZero() {
		updatedBefore = time.Now().UTC()
	}
	institutionID, _ := strconv.ParseInt(c.Query("institution_id"), 10, 64)
	storageOption := c.Query("storage_option")
	chartMetric := c.Query("chart_metric")
	return DepositReportParams{
		ChartMetric:   chartMetric,
		InstitutionID: institutionID,
		StorageOption: storageOption,
		UpdatedBefore: updatedBefore,
	}
}
