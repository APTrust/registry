package main

import (
	"html/template"

	c "github.com/APTrust/registry/controllers"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Set up our template helper functions.
	// These have to be defined before views
	// are loaded, or the view parser will
	// error out.
	r.SetFuncMap(template.FuncMap{
		"dateISO":  helpers.DateISO,
		"dateUS":   helpers.DateUS,
		"truncate": helpers.Truncate,
	})

	// Load the view templates
	r.LoadHTMLGlob("views/**/*.html")
	r.Use(middleware.Auth())
	initRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// initRoutes maps URLs to handlers.
func initRoutes(router *gin.Engine) {

	// This ensures that routes match with or without trailing slash.
	router.RedirectTrailingSlash = true

	router.Static("/static", "./static")
	router.Static("/favicon.ico", "./static/img/favicon.png")

	web := router.Group("/")
	{
		// Dashboard
		web.GET("/dashboard", c.DashboardShow)

		// Users
		web.POST("/users/new", c.UserCreate)
		web.DELETE("/users/delete/:id", c.UserDelete)
		web.GET("/users", c.UserIndex)
		web.GET("/users/new", c.UserNew)
		web.GET("/users/sign_in", c.UserSignInShow)
		web.POST("/users/sign_in", c.UserSignIn)
		web.GET("/users/sign_out", c.UserSignOut) // should be delete?
		web.GET("/users/show/:id", c.UserShow)
		web.PUT("/users/update/:id", c.UserUpdate)

	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "users/sign_in.html", nil)
	})

}
