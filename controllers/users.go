package controllers

import (
	"fmt"
	"net/http"

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

// UserIndex shows list of users.
// GET /users
func UserIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "users/index.html", gin.H{})
}

// UserNew returns a blank form for the user to create a new user.
// GET /users/new
func UserNew(c *gin.Context) {

}

// UserSignIn shows the user sign-in form.
// GET /users/sign_in
func UserSignInShow(c *gin.Context) {
	c.HTML(200, "users/sign_in.html", nil)
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
		})
	}
}

// UserSignOut signs the user out.
// GET /users/sign_out
func UserSignOut(c *gin.Context) {
	helpers.DeleteSessionCookie(c)
	c.HTML(http.StatusOK, "users/sign_in.html", gin.H{
		"error": "You have successfully signed out.",
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
		fmt.Println(err)
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
