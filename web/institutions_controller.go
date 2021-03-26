package web

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// InstitutionCreate a new institution. Handles submission of new
// institution form.
// POST /institutions/new
func InstitutionCreate(c *gin.Context) {
	saveInstitutionForm(c)
}

// InstitutionDelete deletes a user.
// DELETE /institutions/delete/:id
func InstitutionDelete(c *gin.Context) {
	req := NewRequest(c)
	inst, err := pgmodels.InstitutionByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	err = inst.Delete()
	if AbortIfError(c, err) {
		return
	}
	c.Redirect(http.StatusFound, "/institutions")
}

// InstitutionUndelete reactivates an institution.
// GET /institutions/undelete/:id
func InstitutionUndelete(c *gin.Context) {
	req := NewRequest(c)
	inst, err := pgmodels.InstitutionByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	err = inst.Undelete()
	if AbortIfError(c, err) {
		return
	}
	location := fmt.Sprintf("/institutions/show/%d", inst.ID)
	c.Redirect(http.StatusFound, location)
}

// InstitutionIndex shows list of institutions.
// GET /institutions
func InstitutionIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "institutions/index.html"
	query := pgmodels.NewQuery().OrderBy("name")
	institutions, err := pgmodels.InstitutionViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["institutions"] = institutions
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// InstitutionNew returns a blank form for the institution to create
// a new institution.
// GET /institutions/new
func InstitutionNew(c *gin.Context) {
	req := NewRequest(c)
	template := "institutions/form.html"
	form, err := NewInstitutionForm(req)
	if AbortIfError(c, err) {
		return
	}
	form.Action = "/institutions/new"
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// InstitutionShow returns the institution with the specified id.
// GET /institutions/show/:id
func InstitutionShow(c *gin.Context) {
	req := NewRequest(c)
	institution, err := pgmodels.InstitutionViewByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["institution"] = institution

	query := pgmodels.NewQuery().Where("parent_id", "=", institution.ID).OrderBy("name")
	subscribers, err := pgmodels.InstitutionViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["subscribers"] = subscribers

	query = pgmodels.NewQuery().Where("institution_id", "=", institution.ID).IsNull("deactivated_at").OrderBy("name")
	users, err := pgmodels.UserViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["users"] = users

	req.TemplateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "institutions/show.html", req.TemplateData)
}

// InstitutionUpdate saves changes to an exiting institution.
// PUT /institutions/edit/:id
func InstitutionUpdate(c *gin.Context) {
	saveInstitutionForm(c)
}

// InstitutionEdit shows a form to edit an exiting institution.
// GET /institutions/edit/:id
func InstitutionEdit(c *gin.Context) {
	req := NewRequest(c)
	form, err := NewInstitutionForm(req)
	if AbortIfError(c, err) {
		return
	}
	form.Action = fmt.Sprintf("/institutions/edit/%d", form.Model.GetID())
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, "institutions/form.html", req.TemplateData)
}

// TODO: Move common code into Form.
func saveInstitutionForm(c *gin.Context) {
	req := NewRequest(c)
	form, err := NewInstitutionForm(req)
	if AbortIfError(c, err) {
		return
	}

	template := "institutions/form.html"
	form.Action = "/institutions/new"
	if req.ResourceID > 0 {
		form.Action = fmt.Sprintf("/institutions/edit/%d", req.ResourceID)
	}

	req.TemplateData["form"] = form
	status, err := form.Save()
	if err != nil {
		c.HTML(status, template, req.TemplateData)
		return
	}
	location := fmt.Sprintf("/institutions/show/%d?flash=Institution+saved", form.Model.GetID())
	c.Redirect(http.StatusSeeOther, location)
}
