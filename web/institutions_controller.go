package web

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
	saveInstitutionForm(c)
}

// InstitutionDelete deletes a user.
// DELETE /institutions/delete/:id
// GET /institutions/delete/:id
func InstitutionDelete(c *gin.Context) {
	req := NewRequest(c)
	inst, err := pgmodels.InstitutionByID(req.Auth.ResourceID)
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
	inst, err := pgmodels.InstitutionByID(req.Auth.ResourceID)
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
	var institutions []*pgmodels.InstitutionView
	err := req.LoadResourceList(&institutions, "name", forms.NewInstitutionFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// InstitutionNew returns a blank form for the institution to create
// a new institution.
// GET /institutions/new
func InstitutionNew(c *gin.Context) {
	req := NewRequest(c)
	form, err := forms.NewInstitutionForm(&pgmodels.Institution{})
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, form.Template, req.TemplateData)
}

// InstitutionShow returns the institution with the specified id.
// GET /institutions/show/:id
func InstitutionShow(c *gin.Context) {
	req := NewRequest(c)
	institution, err := pgmodels.InstitutionViewByID(req.Auth.ResourceID)
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
	institution, err := pgmodels.InstitutionByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	form, err := forms.NewInstitutionForm(institution)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, form.Template, req.TemplateData)
}

func saveInstitutionForm(c *gin.Context) {
	req := NewRequest(c)
	var err error
	institution := &pgmodels.Institution{}
	if req.Auth.ResourceID > 0 {
		institution, err = pgmodels.InstitutionByID(req.Auth.ResourceID)
		if AbortIfError(c, err) {
			return
		}
	}
	// Bind submitted form values in case we have to
	// re-display the form with an error message.
	c.ShouldBind(institution)
	form, err := forms.NewInstitutionForm(institution)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["form"] = form
	if form.Save() {
		c.Redirect(form.Status, form.PostSaveURL())
	} else {
		req.TemplateData["FormError"] = form.Error
		c.HTML(form.Status, form.Template, req.TemplateData)
	}
}
