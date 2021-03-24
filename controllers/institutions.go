package controllers

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// InstitutionCreate a new institution. Handles submission of new
// institution form.
// POST /institutions/new
func InstitutionCreate(c *gin.Context) {
	saveInstitutionForm(c, 0)
}

// InstitutionDelete deletes a institution.
// DELETE /institutions/:id
func InstitutionDelete(c *gin.Context) {

}

// InstitutionIndex shows list of institutions.
// GET /institutions
func InstitutionIndex(c *gin.Context) {
	r := NewRequest(c)
	template := "institutions/index.html"
	query := pgmodels.NewQuery().OrderBy("name")
	institutions, err := pgmodels.InstitutionViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	r.TemplateData["institutions"] = institutions
	c.HTML(http.StatusOK, template, r.TemplateData)
}

// InstitutionNew returns a blank form for the institution to create
// a new institution.
// GET /institutions/new
func InstitutionNew(c *gin.Context) {
	r := NewRequest(c)
	template := "institutions/form.html"
	form, err := forms.NewInstitutionForm(c, r.TemplateData, 0)
	if AbortIfError(c, err) {
		return
	}
	form.Action = "/institutions/new"
	r.TemplateData["form"] = form
	c.HTML(http.StatusOK, template, r.TemplateData)
}

// InstitutionShow returns the institution with the specified id.
// GET /institutions/show/:id
func InstitutionShow(c *gin.Context) {
	r := NewRequest(c)
	institution, err := pgmodels.InstitutionViewByID(r.ID)
	if AbortIfError(c, err) {
		return
	}
	r.TemplateData["institution"] = institution

	query := pgmodels.NewQuery().Where("parent_id", "=", institution.ID).OrderBy("name")
	subscribers, err := pgmodels.InstitutionViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	r.TemplateData["subscribers"] = subscribers

	query = pgmodels.NewQuery().Where("institution_id", "=", institution.ID).IsNull("deactivated_at").OrderBy("name")
	users, err := pgmodels.UserViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	r.TemplateData["users"] = users

	r.TemplateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "institutions/show.html", r.TemplateData)
}

// InstitutionUpdate saves changes to an exiting institution.
// PUT /institutions/edit/:id
func InstitutionUpdate(c *gin.Context) {
	r := NewRequest(c)
	// institution, err := pgmodels.InstitutionByID(r.ID)
	// if AbortIfError(c, err) {
	// 	return
	// }
	saveInstitutionForm(c, r.ID)
}

// InstitutionEdit shows a form to edit an exiting institution.
// GET /institutions/edit/:id
func InstitutionEdit(c *gin.Context) {
	r := NewRequest(c)
	// institution, err := pgmodels.InstitutionByID(r.ID)
	// if AbortIfError(c, err) {
	// 	return
	// }
	form, err := forms.NewInstitutionForm(c, r.TemplateData, r.ID)
	if AbortIfError(c, err) {
		return
	}
	form.Action = fmt.Sprintf("/institutions/edit/%d", form.Model.GetID())
	r.TemplateData["form"] = form
	c.HTML(http.StatusOK, "institutions/form.html", r.TemplateData)
}

func saveInstitutionForm(c *gin.Context, instID int64) {
	r := NewRequest(c)
	form, err := forms.NewInstitutionForm(c, r.TemplateData, instID)
	if AbortIfError(c, err) {
		return
	}

	template := "institutions/form.html"
	form.Action = "/institutions/new"
	if instID > 0 {
		form.Action = fmt.Sprintf("/institutions/edit/%d", instID)
	}

	r.TemplateData["form"] = form
	status, err := form.Save()
	if err != nil {
		c.HTML(status, template, r.TemplateData)
		return
	}
	location := fmt.Sprintf("/institutions/show/%d?flash=Institution+saved", form.Model.GetID())
	c.Redirect(http.StatusSeeOther, location)
}
