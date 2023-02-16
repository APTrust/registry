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
	ReportType    string
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
	var deposits []*pgmodels.DepositStats
	var err error
	if params.ReportType == "over_time" {
		deposits, err = pgmodels.DepositStatsOverTime(params.InstitutionID, params.StorageOption)
	} else {
		deposits, err = pgmodels.DepositStatsSelect(params.InstitutionID, params.StorageOption, params.UpdatedBefore)
	}
	if AbortIfError(c, err) {
		return
	}
	filterCollection := req.GetFilterCollection()
	if filterCollection.ValueOf("report_type") == "" {
		filterCollection.Add("report_type", []string{params.ReportType})
	}
	filterForm, err := forms.NewDepositReportFilterForm(filterCollection, req.CurrentUser)
	if AbortIfError(c, err) {
		return
	}

	if params.ReportType == "over_time" {
		fields := filterForm.GetFields()
		fields["storage_option"].Attrs["disabled"] = "true"
		fields["end_date"].Attrs["disabled"] = "true"

		// Time report covers through end of prior month.
		if len(fields["end_date"].Options) > 1 {
			fields["end_date"].Value = fields["end_date"].Options[1].Value
		}
	}

	instList := depositInstList(deposits)
	storageOptionsList := depositStorageOptions(deposits)

	req.TemplateData["reportType"] = params.ReportType
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
	updatedBefore, _ := time.Parse("2006-01-02", c.Query("end_date"))
	if updatedBefore.IsZero() {
		updatedBefore = time.Now().UTC()
	}
	institutionID, _ := strconv.ParseInt(c.Query("institution_id"), 10, 64)
	storageOption := c.Query("storage_option")
	chartMetric := c.Query("chart_metric")
	reportType := c.Query("report_type")
	if reportType == "" {
		reportType = "by_inst"
	}
	return DepositReportParams{
		ChartMetric:   chartMetric,
		InstitutionID: institutionID,
		ReportType:    reportType,
		StorageOption: storageOption,
		UpdatedBefore: updatedBefore,
	}
}
