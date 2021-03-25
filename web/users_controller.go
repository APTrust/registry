package web

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserCreate a new user. Handles submission of new user form.
// POST /users/new
func UserCreate(c *gin.Context) {
	saveUserForm(c)
}

// UserDelete deletes a user.
// DELETE /users/:id
func UserDelete(c *gin.Context) {

}

// UserIndex shows list of users.
// GET /users
func UserIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "users/index.html"
	query, err := getIndexQuery(c)
	if AbortIfError(c, err) {
		return
	}

	query.OrderBy("name asc")
	req.TemplateData["selectedID"] = c.Query("institution_id")

	users, err := pgmodels.UserViewSelect(query)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["users"] = users

	// Get institutions for filter list
	institutionOptions, err := ListInstitutions(false)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["institutionOptions"] = institutionOptions

	c.HTML(http.StatusOK, template, req.TemplateData)
}

// UserNew returns a blank form for the user to create a new user.
// GET /users/new
func UserNew(c *gin.Context) {
	req := NewRequest(c)
	template := "users/form.html"
	form, err := NewUserForm(req)
	if AbortIfError(c, err) {
		return
	}
	form.Action = "/users/new"
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// UserShow returns the user with the specified id.
// GET /users/show/:id
func UserShow(c *gin.Context) {
	req := NewRequest(c)
	user, err := pgmodels.UserByID(req.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["user"] = user
	req.TemplateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "users/show.html", req.TemplateData)
}

// UserUpdate saves changes to an exiting user.
// PUT /users/edit/:id
func UserUpdate(c *gin.Context) {
	saveUserForm(c)
}

// UserEdit shows a form to edit an exiting user.
// GET /users/edit/:id
func UserEdit(c *gin.Context) {
	req := NewRequest(c)
	form, err := NewUserForm(req)
	if AbortIfError(c, err) {
		return
	}
	form.Action = fmt.Sprintf("/users/edit/%d", req.ResourceID)
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, "users/form.html", req.TemplateData)
}

// UserSignInShow shows the user sign-in form.
// GET /users/sign_in
func UserSignInShow(c *gin.Context) {
	c.HTML(200, "users/sign_in.html", gin.H{
		"cover": helpers.GetCover(),
	})
}

// UserSignIn signs the user in.
// POST /users/sign_in
func UserSignIn(c *gin.Context) {
	status, redirectTo, err := SignInUser(c)
	if err == nil {
		c.Redirect(status, redirectTo)
	} else {
		c.HTML(status, "users/sign_in.html", gin.H{
			"error": err.Error(),
			"cover": helpers.GetCover(),
		})
	}
}

// UserSignOut signs the user out.
// GET /users/sign_out
func UserSignOut(c *gin.Context) {
	user := helpers.CurrentUser(c)
	if user != nil {
		user.SignOut()
	}
	helpers.DeleteSessionCookie(c)
	c.HTML(http.StatusOK, "users/sign_in.html", gin.H{
		"cover": helpers.GetCover(),
	})
}

func SignInUser(c *gin.Context) (int, string, error) {
	redirectTo := "/users/sign_in"
	user, err := pgmodels.UserSignIn(
		c.PostForm("email"),
		c.PostForm("password"),
		c.ClientIP(),
	)
	if err != nil {
		c.Error(err)
		helpers.DeleteSessionCookie(c)
		return http.StatusBadRequest, redirectTo, err
	}
	err = helpers.SetSessionCookie(c, user)
	if err != nil {
		return http.StatusInternalServerError, redirectTo, err
	}
	c.Set("CurrentUser", user)
	startPage := "/dashboard"
	if !user.IsAdmin() {
		startPage += fmt.Sprintf("?institution_id=%d", user.InstitutionID)
	}
	return http.StatusFound, startPage, nil
}

func getIndexQuery(c *gin.Context) (*pgmodels.Query, error) {
	allowedFilters := []string{
		"institution_id",
	}
	fc := NewFilterCollection()
	for _, key := range allowedFilters {
		fc.Add(key, c.QueryArray(key))
	}
	return fc.ToQuery()
}

// TODO: Move common code into Form.
func saveUserForm(c *gin.Context) {
	req := NewRequest(c)
	form, err := NewUserForm(req)
	if AbortIfError(c, err) {
		return
	}

	template := "users/form.html"
	form.Action = "/users/new"
	if req.ResourceID > 0 {
		form.Action = fmt.Sprintf("/users/edit/%d", req.ResourceID)
	} else {
		// Assign random password to new user. They'll get an email
		// asking them to reset their password.
		encPwd, err := common.EncryptPassword(uuid.New().String())
		if AbortIfError(c, err) {
			return
		}
		form.Model.(*pgmodels.User).EncryptedPassword = encPwd
	}

	req.TemplateData["form"] = form
	status, err := form.Save()
	if err != nil {
		c.HTML(status, template, req.TemplateData)
		return
	}
	location := fmt.Sprintf("/users/show/%d?flash=User+saved", form.Model.GetID())
	c.Redirect(http.StatusSeeOther, location)
}
