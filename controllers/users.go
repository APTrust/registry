package controllers

import (
	"net/http"
	"strconv"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/models"
	"github.com/gin-gonic/gin"
)

// UserCreate a new user. Handles submission of new user form.
// POST /users/new
func UserCreate(c *gin.Context) {

}

// UserDelete deletes a user.
// DELETE /users/:id
func UserDelete(c *gin.Context) {

}

// TODO: Query users_view instead.
// Allow filtering on name, email, institution, role, deactivated.
// Allow sorting on name, email, institution, created,
// updated, current sign in, last sign in,

// UserIndex shows list of users.
// GET /users
func UserIndex(c *gin.Context) {

	// TODO: Create common convenience methods to get select list
	//       data. For example, this controller uses an Institutions
	//       list, so models.InstitutionsAll?
	//

	// OK if instID doesn't parse. It will often be zero.

	instID, _ := strconv.ParseInt(c.Query("institution_id__eq"), 10, 64)

	template := "users/index.html"
	templateData := helpers.TemplateVars(c)
	query, err := getIndexQuery(c)
	if err != nil {
		c.AbortWithError(StatusCodeForError(err), err)
		return
	}
	currentUser := helpers.CurrentUser(c)
	if currentUser == nil {
		err = common.ErrPermissionDenied
		c.AbortWithError(StatusCodeForError(err), err)
		return
	} else if currentUser.IsAdmin() == false {
		query.Where("institution_id", "=", currentUser.InstitutionID)
	}
	// Get user list
	query.Columns(
		"name",
		"email",
		"institution_name",
		"institution_id",
		"role",
		"enabled_two_factor",
		"deactivated_at",
	)
	templateData["selectedID"] = instID

	ds := models.NewDataStore(currentUser)
	users, err := ds.UserViewList(query)
	if err != nil {
		c.AbortWithError(StatusCodeForError(err), err)
		return
	}
	templateData["users"] = users

	// Get institutions
	instQuery := models.NewQuery().Columns("id", "name").OrderBy("name asc").Limit(100).Offset(0)
	institutions, err := ds.InstitutionList(instQuery)
	if err != nil {
		c.AbortWithError(StatusCodeForError(err), err)
		return
	}
	templateData["institutions"] = institutions

	c.HTML(http.StatusOK, template, templateData)
}

// UserNew returns a blank form for the user to create a new user.
// GET /users/new
func UserNew(c *gin.Context) {

}

// UserSignIn shows the user sign-in form.
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
	c.HTML(http.StatusOK, "users/show.html", gin.H{})
}

// UserUpdate saves changes to an exiting user.
// PUT /users/:id
func UserUpdate(c *gin.Context) {

}

func SignInUser(c *gin.Context) (int, string, error) {
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
