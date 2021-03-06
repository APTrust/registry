package web

import (
	"github.com/APTrust/registry/pgmodels"

	"github.com/gin-gonic/gin"
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
	err = loadDashDeposits(r)
	return err
}

func loadDashWorkItems(r *Request) error {
	query := pgmodels.NewQuery().OrderBy("date_processed desc").Offset(0).Limit(10)
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
	query := pgmodels.NewQuery().Where("user_id", "=", r.Auth.CurrentUser().ID).OrderBy("created_at desc").Offset(0).Limit(10)
	alerts, err := pgmodels.AlertViewSelect(query)
	if err != nil {
		return err
	}
	r.TemplateData["alerts"] = alerts
	return nil
}

func loadDashDeposits(r *Request) error {
	query := pgmodels.NewQuery().OrderBy("institution_name asc").OrderBy("storage_option asc")
	if r.Auth.CurrentUser().IsAdmin() {
		query.IsNull("institution_id")
	} else {
		query.Where("institution_id", "=", r.Auth.CurrentUser().InstitutionID)
	}
	stats, err := pgmodels.StorageOptionStatsSelect(query)
	if err != nil {
		return err
	}
	r.TemplateData["stats"] = stats
	return nil
}
