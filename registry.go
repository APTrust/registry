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
		"escapeAttr": helpers.EscapeAttr,
		"escapeHTML": helpers.EscapeHTML,
		"strEq":      helpers.StrEq,
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
// Note that for PUT requests, we also have to support POST,
// and for DELETE, we need to support GET. This is in the web
// routes only, because browsers don't support PUT and DELETE,
// only POST and GET. Most other frameworks will rewrite the
// request verb based on a _method param in the form or query
// string, but gin won't because the router gets its hands on the
// request before any other middleware. This issue is documented
// in a few GitHub issues, including
// https://github.com/gin-gonic/gin/issues/450
//
// The maintainers' solution is to use r.Any(), meaning any HTTP
// verb would map to a given route. We'll use the pairs PUT/POST
// and GET/DELETE, which is a little more restrictive that "match
// anything".
func initRoutes(router *gin.Engine) {

	// This ensures that routes match with or without trailing slash.
	router.RedirectTrailingSlash = true

	router.Static("/static", "./static")
	router.Static("/favicon.ico", "./static/img/favicon.png")

	web := router.Group("/")
	{
		// Dashboard
		web.GET("/dashboard", c.DashboardShow)

		// Institutions
		web.POST("/institutions/new", c.InstitutionCreate)
		web.DELETE("/institutions/delete/:id", c.InstitutionDelete)
		web.GET("/institutions/delete/:id", c.InstitutionDelete)
		web.GET("/institutions", c.InstitutionIndex)
		web.GET("/institutions/new", c.InstitutionNew)
		web.GET("/institutions/show/:id", c.InstitutionShow)
		web.GET("/institutions/edit/:id", c.InstitutionEdit)
		web.PUT("/institutions/edit/:id", c.InstitutionUpdate)
		web.POST("/institutions/edit/:id", c.InstitutionUpdate)

		// Users
		web.POST("/users/new", c.UserCreate)
		web.DELETE("/users/delete/:id", c.UserDelete)
		web.GET("/users/delete/:id", c.UserDelete)
		web.GET("/users", c.UserIndex)
		web.GET("/users/new", c.UserNew)
		web.GET("/users/show/:id", c.UserShow)
		web.GET("/users/edit/:id", c.UserEdit)
		web.PUT("/users/edit/:id", c.UserUpdate)
		web.POST("/users/edit/:id", c.UserUpdate)

		// User Sign In
		web.GET("/users/sign_in", c.UserSignInShow)
		web.POST("/users/sign_in", c.UserSignIn)
		web.GET("/users/sign_out", c.UserSignOut) // should be delete?

		// Error page
		web.GET("/error", c.ErrorShow)

	}

	// Root goes to sign-in page
	router.GET("/", c.UserSignInShow)
}
