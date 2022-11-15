package webui

import (
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func DashboardShow(c *gin.Context) {
	r := NewRequest(c)
	err := loadDashData(r)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(200, "dashboard/show.html", r.TemplateData)
}

func loadDashData(r *Request) error {
	err := loadDashWorkItems(r)
	if err != nil {
		return err
	}
	err = loadDashAlerts(r)
	if err != nil {
		return err
	}
	loadDashStats(r)
	err = loadDashDeposits(r)
	return err
}

func loadDashWorkItems(r *Request) error {
	query := pgmodels.NewQuery().OrderBy("date_processed", "desc").Offset(0).Limit(10)
	if !r.Auth.CurrentUser().IsAdmin() {
		query.Where("institution_id", "=", r.Auth.CurrentUser().InstitutionID)
	}
	items, err := pgmodels.WorkItemViewSelect(query)
	if err != nil {
		return err
	}
	r.TemplateData["items"] = items
	return nil
}

func loadDashAlerts(r *Request) error {
	query := pgmodels.NewQuery().Where("user_id", "=", r.Auth.CurrentUser().ID).OrderBy("created_at", "desc").Offset(0).Limit(10)
	alerts, err := pgmodels.AlertViewSelect(query)
	if err != nil {
		return err
	}
	r.TemplateData["alerts"] = alerts
	return nil
}

func loadDashDeposits(r *Request) error {
	institutionID := r.Auth.CurrentUser().InstitutionID
	if r.Auth.CurrentUser().IsAdmin() {
		institutionID = 0
	}
	stats, err := pgmodels.DepositStatsSelect(institutionID, "", time.Now().UTC())
	if err != nil {
		return err
	}

	// This is a bit of a hack. We should add this to the query,
	// but the query is already complex...
	filteredStats := make([]*pgmodels.DepositStats, 0)
	for _, stat := range stats {
		if stat.InstitutionID == institutionID {
			filteredStats = append(filteredStats, stat)
		}
	}

	r.TemplateData["depositStats"] = filteredStats
	return nil
}

func loadDashStats(r *Request) {

	p := message.NewPrinter(language.English)

	query := pgmodels.NewQuery()
	if !r.Auth.CurrentUser().IsAdmin() {
		query.Where("institution_id", "=", r.Auth.CurrentUser().InstitutionID)
	}

	var objs []*pgmodels.IntellectualObjectView
	objCount, err := pgmodels.GetCountFromView(query, objs)
	if err != nil {
		common.Context().Log.Warn().Msgf("error running object count query for dashboard: %v", err)
	}
	r.TemplateData["objectCount"] = p.Sprintf("%d", objCount)

	// For objects and files, we want to count only Active items
	query.Where("state", "=", "A")

	var files []*pgmodels.GenericFileView
	fileCount, err := pgmodels.GetCountFromView(query, files)
	if err != nil {
		common.Context().Log.Warn().Msgf("error running file count query for dashboard: %v", err)
	}
	r.TemplateData["fileCount"] = p.Sprintf("%d", fileCount)

	var events []*pgmodels.PremisEventView
	eventCount, err := pgmodels.GetCountFromView(query, events)
	if err != nil {
		common.Context().Log.Warn().Msgf("error running premis event count query for dashboard: %v", err)
	}
	r.TemplateData["eventCount"] = p.Sprintf("%d", eventCount)

}
