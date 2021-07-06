package web

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/forms"
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
// DELETE /users/delete/:id
func UserDelete(c *gin.Context) {
	req := NewRequest(c)
	user, err := pgmodels.UserByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	err = user.Delete()
	if AbortIfError(c, err) {
		return
	}
	location := fmt.Sprintf("/users?institution_id=%d", user.InstitutionID)
	c.Redirect(http.StatusFound, location)
}

// UserUndelete reactivates a user.
// GET /users/undelete/:id
func UserUndelete(c *gin.Context) {
	req := NewRequest(c)
	user, err := pgmodels.UserByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	err = user.Undelete()
	if AbortIfError(c, err) {
		return
	}
	location := fmt.Sprintf("/users/show/%d", user.ID)
	c.Redirect(http.StatusFound, location)
}

// UserIndex shows list of users.
// GET /users
func UserIndex(c *gin.Context) {
	req := NewRequest(c)
	template := "users/index.html"
	var users []*pgmodels.UserView
	err := req.LoadResourceList(&users, "name asc", forms.NewUserFilterForm)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, template, req.TemplateData)
}

// UserNew returns a blank form for the user to create a new user.
// GET /users/new
func UserNew(c *gin.Context) {
	req := NewRequest(c)
	form, err := forms.NewUserForm(&pgmodels.User{}, req.CurrentUser)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, form.Template, req.TemplateData)
}

// UserShow returns the user with the specified id.
// GET /users/show/:id
func UserShow(c *gin.Context) {
	req := NewRequest(c)
	user, err := pgmodels.UserByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["user"] = user
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
	userToEdit, err := pgmodels.UserByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	form, err := forms.NewUserForm(userToEdit, req.CurrentUser)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, form.Template, req.TemplateData)
}

// UserSignInShow shows the user sign-in form.
// GET /users/sign_in
func UserSignInShow(c *gin.Context) {
	c.HTML(200, "users/sign_in.html", gin.H{
		"cover":             helpers.GetCover(),
		"preFillTestLogins": common.Context().Config.EnvName == "test",
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
		"cover":             helpers.GetCover(),
		"preFillTestLogins": common.Context().Config.EnvName == "test",
	})
}

// UserShowChangePassword displays the change password page
// for the user with the specified ID.
//
// GET /users/change_password/:id
func UserShowChangePassword(c *gin.Context) {
	req := NewRequest(c)
	if req.CurrentUser.ID != req.Auth.ResourceID && !req.CurrentUser.IsAdmin() {
		AbortIfError(c, common.ErrPermissionDenied)
		return
	}
	userToEdit, err := pgmodels.UserByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	form := forms.NewPasswordResetForm(userToEdit)
	req.TemplateData["form"] = form
	c.HTML(http.StatusOK, form.Template, req.TemplateData)
}

// UserChangePassword changes a user's password. The user gets
// to specify what the new password will be.
//
// POST /users/change_password/:id
func UserChangePassword(c *gin.Context) {
	// Force error if CurrentUser.ID != req.Auth.ResourceID
}

// UserForcePasswordReset resets a user's password to something
// random and sends them an email with a link that has an embedded
// login token. When they follow the link, they'll be automatically
// logged in and will have to choose a new password.
//
// POST /users/force_password_reset/:id
func UserForcePasswordReset(c *gin.Context) {
	// This is admin triggering a password reset for another user,
	// so current user id does not need to match subject user id.
}

// UserGetAPIKey issues a new API key for the user, which replaces the
// user's existing API key. This key will be displayed once only.
//
// POST /users/get_api_key/:id
func UserGetAPIKey(c *gin.Context) {
	// Force error if CurrentUser.ID != req.Auth.ResourceID
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
	fc := pgmodels.NewFilterCollection()
	for _, key := range allowedFilters {
		fc.Add(key, c.QueryArray(key))
	}
	return fc.ToQuery()
}

func saveUserForm(c *gin.Context) {
	req := NewRequest(c)
	userToEdit := &pgmodels.User{}
	var err error
	if req.Auth.ResourceID > 0 {
		// Load existing user.
		userToEdit, err = pgmodels.UserByID(req.Auth.ResourceID)
		if AbortIfError(c, err) {
			return
		}
	} else {
		// Assign random password to new user. They'll get an email
		// asking them to reset their password.
		encPwd, err := common.EncryptPassword(uuid.New().String())
		if AbortIfError(c, err) {
			return
		}
		userToEdit.EncryptedPassword = encPwd
	}

	// Bind submitted form values in case we have to
	// re-display the form with an error message.
	c.ShouldBind(userToEdit)
	form, err := forms.NewUserForm(userToEdit, req.CurrentUser)
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
