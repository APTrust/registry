package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
)

// UserCreate a new user. Handles submission of new user form.
// POST /users/new
func UserCreate(c *gin.Context) {
	currentUser := helpers.CurrentUser(c)
	ds := models.NewDataStore(currentUser)
	template := "users/form.html"
	templateData := helpers.TemplateVars(c)

	form, err := forms.NewUserForm(ds, &models.User{})
	if AbortIfError(c, err) {
		return
	}
	form.Action = "/users/new"
	templateData["form"] = form
	err = form.Bind(c)
	// If validation error, re-display the form with error messages.
	if err != nil {
		c.HTML(http.StatusBadRequest, template, templateData)
		return
	}

	// If no validation error, save the user and redirect.
	err = ds.UserSave(form.User)
	if AbortIfError(c, err) {
		return
	}
	location := fmt.Sprintf("/users/show/%d?flash=User+created", form.User.ID)
	c.Redirect(http.StatusSeeOther, location)
}

// UserDelete deletes a user.
// DELETE /users/:id
func UserDelete(c *gin.Context) {

}

// UserIndex shows list of users.
// GET /users
func UserIndex(c *gin.Context) {
	template := "users/index.html"
	templateData := helpers.TemplateVars(c)
	query, err := getIndexQuery(c)
	if AbortIfError(c, err) {
		return
	}

	query.OrderBy("name asc")
	currentUser := helpers.CurrentUser(c)
	templateData["selectedID"] = c.Query("institution_id__eq")

	ds := models.NewDataStore(currentUser)
	users, err := ds.UserViewList(query)
	if AbortIfError(c, err) {
		return
	}
	templateData["users"] = users

	// Get institutions for filter list
	institutionOptions, err := forms.ListInstitutions(ds)
	if AbortIfError(c, err) {
		return
	}
	templateData["institutionOptions"] = institutionOptions

	c.HTML(http.StatusOK, template, templateData)
}

// UserNew returns a blank form for the user to create a new user.
// GET /users/new
func UserNew(c *gin.Context) {
	currentUser := helpers.CurrentUser(c)
	ds := models.NewDataStore(currentUser)
	template := "users/form.html"
	templateData := helpers.TemplateVars(c)
	form, err := forms.NewUserForm(ds, &models.User{})
	if AbortIfError(c, err) {
		return
	}
	form.Action = "/users/new"
	templateData["form"] = form
	c.HTML(http.StatusOK, template, templateData)
}

// UserShow returns the user with the specified id.
// GET /users/show/:id
func UserShow(c *gin.Context) {
	templateData := helpers.TemplateVars(c)
	currentUser := helpers.CurrentUser(c)
	ds := models.NewDataStore(currentUser)
	user, err := findUser(ds, c.Param("id"))
	if AbortIfError(c, err) {
		return
	}
	templateData["user"] = user
	templateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "users/show.html", templateData)
}

// UserUpdate saves changes to an exiting user.
// PUT /users/edit/:id
func UserUpdate(c *gin.Context) {

}

// UserEdit shows a form to edit an exiting user.
// GET /users/edit/:id
func UserEdit(c *gin.Context) {
	templateData := helpers.TemplateVars(c)
	currentUser := helpers.CurrentUser(c)
	ds := models.NewDataStore(currentUser)
	userToEdit, err := findUser(ds, c.Param("id"))
	if AbortIfError(c, err) {
		return
	}
	form, err := forms.NewUserForm(ds, userToEdit)
	if AbortIfError(c, err) {
		return
	}
	form.Action = fmt.Sprintf("/users/edit/%d", userToEdit.ID)
	templateData["form"] = form
	c.HTML(http.StatusOK, "users/form.html", templateData)
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
	helpers.DeleteSessionCookie(c)
	c.HTML(http.StatusOK, "users/sign_in.html", gin.H{
		"cover": helpers.GetCover(),
	})
}

func SignInUser(c *gin.Context) (int, string, error) {
	// Second of two DataStore instances with automatic
	// admin privileges.
	ds := models.NewDataStore(&models.User{Role: constants.RoleSysAdmin})
	redirectTo := "/users/sign_in"
	user, err := ds.UserSignIn(
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
	return http.StatusFound, "/dashboard", nil
}

func findUser(ds *models.DataStore, id string) (*models.User, error) {
	userID, _ := strconv.ParseInt(id, 10, 64)
	return ds.UserFind(userID)
}

func getIndexQuery(c *gin.Context) (*models.Query, error) {
	allowedFilters := []string{
		"institution_id__eq",
	}
	fc := NewFilterCollection()
	for _, key := range allowedFilters {
		fc.Add(key, c.QueryArray(key))
	}
	return fc.ToQuery()
}
