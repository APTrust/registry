package main

import (
	"html/template"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	"github.com/APTrust/registry/web"
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
		"dateTimeISO":    helpers.DateTimeISO,
		"dateUS":         helpers.DateUS,
		"dict":           helpers.Dict,
		"escapeAttr":     helpers.EscapeAttr,
		"escapeHTML":     helpers.EscapeHTML,
		"humanSize":      helpers.HumanSize,
		"iconFor":        helpers.IconFor,
		"replace":        strings.Replace,
		"roleName":       helpers.RoleName,
		"strEq":          helpers.StrEq,
		"titleCase":      strings.Title,
		"truncate":       helpers.Truncate,
		"truncateMiddle": helpers.TruncateMiddle,
		"truncateStart":  helpers.TruncateStart,
		"userCan":        helpers.UserCan,
		"yesNo":          helpers.YesNo,
		"dateISO":        helpers.DateISO,
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

	// Then authentication and authorization middleware
	router.Use(middleware.Authenticate())
	router.Use(middleware.Authorize())
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
// TODO: Implement the hidden _method hack described here, if possible.
// https://stackoverflow.com/questions/16805956/why-dont-browsers-support-put-and-delete-requests-and-when-will-they
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

	webRoutes := router.Group("/")
	{
		// Dashboard
		webRoutes.GET("/dashboard", web.DashboardShow)

		// Institutions
		webRoutes.POST("/institutions/new", web.InstitutionCreate)
		webRoutes.DELETE("/institutions/delete/:id", web.InstitutionDelete)
		webRoutes.GET("/institutions/delete/:id", web.InstitutionDelete)
		webRoutes.GET("/institutions/undelete/:id", web.InstitutionUndelete)
		webRoutes.GET("/institutions", web.InstitutionIndex)
		webRoutes.GET("/institutions/new", web.InstitutionNew)
		webRoutes.GET("/institutions/show/:id", web.InstitutionShow)
		webRoutes.GET("/institutions/edit/:id", web.InstitutionEdit)
		webRoutes.PUT("/institutions/edit/:id", web.InstitutionUpdate)
		webRoutes.POST("/institutions/edit/:id", web.InstitutionUpdate)

		// GenericFiles
		webRoutes.GET("/files", web.GenericFileIndex)
		webRoutes.GET("/files/show/:id", web.GenericFileShow)
		webRoutes.GET("/files/delete/:id", web.GenericFileDelete)
		webRoutes.GET("/files/restore/:id", web.GenericFileRestore)

		// IntellectualObjects
		webRoutes.GET("/objects", web.IntellectualObjectIndex)
		webRoutes.GET("/objects/show/:id", web.IntellectualObjectShow)
		webRoutes.GET("/objects/request_delete/:id", web.IntellectualObjectRequestDelete)
		webRoutes.DELETE("/objects/init_delete/:id", web.IntellectualObjectInitDelete)
		webRoutes.GET("/objects/restore/:id", web.IntellectualObjectRestore)
		webRoutes.GET("/objects/request_restore/:id", web.IntellectualObjectRequestRestore)
		webRoutes.POST("/objects/init_restore/:id", web.IntellectualObjectInitRestore)
		webRoutes.GET("/objects/events/:id", web.IntellectualObjectEvents)

		// WorkItems - Web UI allows only list, show, and limited editing for admin only
		webRoutes.GET("/work_items", web.WorkItemIndex)
		webRoutes.GET("/work_items/show/:id", web.WorkItemShow)
		webRoutes.GET("/work_items/edit/:id", web.WorkItemEdit)
		webRoutes.PUT("/work_items/edit/:id", web.WorkItemUpdate)
		webRoutes.POST("/work_items/edit/:id", web.WorkItemUpdate)
		webRoutes.PUT("/work_items/requeue/:id", web.WorkItemRequeue)
		webRoutes.POST("/work_items/requeue/:id", web.WorkItemRequeue)

		// Users
		webRoutes.POST("/users/new", web.UserCreate)
		webRoutes.DELETE("/users/delete/:id", web.UserDelete)
		webRoutes.GET("/users/delete/:id", web.UserDelete)
		webRoutes.GET("/users/undelete/:id", web.UserUndelete)
		webRoutes.GET("/users", web.UserIndex)
		webRoutes.GET("/users/new", web.UserNew)
		webRoutes.GET("/users/show/:id", web.UserShow)
		webRoutes.GET("/users/edit/:id", web.UserEdit)
		webRoutes.PUT("/users/edit/:id", web.UserUpdate)
		webRoutes.POST("/users/edit/:id", web.UserUpdate)

		// User Sign In
		webRoutes.GET("/users/sign_in", web.UserSignInShow)
		webRoutes.POST("/users/sign_in", web.UserSignIn)
		webRoutes.GET("/users/sign_out", web.UserSignOut) // should be delete?

		// Error page
		webRoutes.GET("/error", web.ErrorShow)

	}

	// Root goes to sign-in page
	router.GET("/", web.UserSignInShow)
}
