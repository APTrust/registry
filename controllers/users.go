package controllers

import (
	"net/http"

	"github.com/APTrust/registry/common"
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
	allowedFilters := []string{
		"institution_id",
	}
	ctx := common.Context()
	resp := NewIndexRequest(c, allowedFilters, true, "users/index.html")

	userFilter := &models.User{}
	err := c.ShouldBindQuery(userFilter)
	if err != nil {
		resp.SetError(err)
		resp.Respond()
		return
	}
	currentUser := helpers.CurrentUser(c)
	if currentUser == nil {
		resp.SetError(common.ErrPermissionDenied)
	} else if currentUser.IsAdmin() == false {
		userFilter.InstitutionID = currentUser.InstitutionID
	}

	// Get users
	users := make([]*models.UsersView, 0)
	userQuery := ctx.DB.Model(&users).Column("name", "email", "institution_name", "institution_id", "role", "enabled_two_factor", "deactivated_at").Order("name asc")
	if userFilter.InstitutionID > 0 {
		userQuery = userQuery.Where("institution_id = ?", userFilter.InstitutionID)
		resp.TemplateData["selectedID"] = userFilter.InstitutionID
	} else {
		resp.TemplateData["selectedID"] = 0
	}
	resp.SetError(userQuery.Select())
	resp.TemplateData["users"] = users

	// Get institutions
	institutions := make([]*models.Institution, 0)
	instQuery := ctx.DB.Model(&institutions).Column("id", "name").Order("name asc")
	resp.SetError(instQuery.Select())
	resp.TemplateData["institutions"] = institutions

	resp.Respond()
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
	redirectTo := "/users/sign_in"
	user, err := models.SignInUser(
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
