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
	user := &models.User{}
	templateData["user"] = user
	templateData["roles"] = forms.RolesList
	institutionOptions, err := forms.ListInstitutions(ds)
	if AbortIfError(c, err) {
		return
	}
	templateData["institutionOptions"] = institutionOptions
	if err := c.ShouldBind(user); err != nil {
		errMessages := ValidationErrors(err, user)
		if errMessages != nil {
			templateData["validationErrors"] = errMessages
		}
		templateData["error"] = err.Error()
		c.HTML(http.StatusBadRequest, template, templateData)
		return
	}
	err = ds.UserSave(user)
	if AbortIfError(c, err) {
		return
	}
	location := fmt.Sprintf("/users/show/%d?flash=User+created", user.ID)
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
	templateData["user"] = &models.User{}
	templateData["roles"] = forms.RolesList
	institutionOptions, err := forms.ListInstitutions(ds)
	if AbortIfError(c, err) {
		return
	}
	templateData["institutionOptions"] = institutionOptions
	c.HTML(http.StatusOK, template, templateData)
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

// UserShow returns the user with the specified id.
// GET /users/show/:id
func UserShow(c *gin.Context) {
	templateData := helpers.TemplateVars(c)
	currentUser := helpers.CurrentUser(c)
	ds := models.NewDataStore(currentUser)
	id := c.Param("id")
	userID, _ := strconv.ParseInt(id, 10, 64)
	user, err := ds.UserFind(userID)
	if AbortIfError(c, err) {
		return
	}
	templateData["user"] = user
	templateData["flash"] = c.Query("flash")
	c.HTML(http.StatusOK, "users/show.html", templateData)
}

// UserUpdate saves changes to an exiting user.
// PUT /users/:id
func UserUpdate(c *gin.Context) {

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
