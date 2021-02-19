package main

import (
	"html/template"

	"github.com/APTrust/registry/common"
	c "github.com/APTrust/registry/controllers"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	initTemplates(r)
	initMiddleware(r)
	initRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// initTemplateHelpers sets up our template helper functions.
// These have to be defined before views  are loaded, or the view
// parser will error out.
func initTemplates(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"dateISO":    helpers.DateISO,
		"dateUS":     helpers.DateUS,
		"eqStrInt":   helpers.EqStrInt,
		"eqStrInt64": helpers.EqStrInt64,
		"truncate":   helpers.Truncate,
		"roleName":   helpers.RoleName,
		"yesNo":      helpers.YesNo,
	})
	// Load the view templates
	router.LoadHTMLGlob("views/**/*.html")
}

// initMiddleware loads our custom middleware in the desired order.
func initMiddleware(router *gin.Engine) {
	// Logger first...
	ctx := common.Context()
	router.Use(logger.SetLogger(logger.Config{
		Logger: &ctx.Log,
		UTC:    true,
	}))

	// Then authentication middleware
	router.Use(middleware.Auth())
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

		// Error page
		web.GET("/error", c.ErrorShow)

	}

	// Root goes to sign-in page
	router.GET("/", c.UserSignInShow)
}
