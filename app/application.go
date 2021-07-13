package app

import (
	"html/template"
	"io"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	"github.com/APTrust/registry/web"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

// Run runs the Registry application. This is called from main() to start
// the app.
func Run() {
	r := InitAppEngine(false)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// InitAppEngine sets up the whole Gin application, loading templates and
// middleware and defining routes. The test suite can use this to get an
// instance of the Gin engine to bind to.
//
// Set param discardStdOut during unit/integration tests to suppress
// Gin's STDOUT logging. Those log statements are useful in development,
// but can be verbose and clutter the test output.
func InitAppEngine(discardStdOut bool) *gin.Engine {
	var r *gin.Engine
	if discardStdOut {
		r = gin.New()
		r.Use(gin.Recovery())
		gin.DefaultWriter = io.Discard
	} else {
		r = gin.Default()
	}
	initTemplates(r)
	initMiddleware(r)
	initRoutes(r)
	return r
}

// initTemplateHelpers sets up our template helper functions.
// These have to be defined before views  are loaded, or the view
// parser will error out.
func initTemplates(router *gin.Engine) {
	router.SetFuncMap(template.FuncMap{
		"dateTimeISO":    helpers.DateTimeISO,
		"dateUS":         helpers.DateUS,
		"defaultString":  helpers.DefaultString,
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
	// If we're running from main, templates will come
	// from ./views. When running tests, templates come
	// from ../views because http tests run from web
	// sub directory.
	if common.FileExists("./views") {
		router.LoadHTMLGlob("./views/**/*.html")
	} else {
		router.LoadHTMLGlob("../views/**/*.html")
	}
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
		// Alerts
		webRoutes.GET("/alerts/", web.AlertIndex)
		webRoutes.GET("/alerts/show/:id/:user_id", web.AlertShow)

		// Deletion Requests
		// Note that these routes are for read-only views.
		// Routes for initiating, approving and rejecting deletions
		// are in the GenericFiles and IntellectualObjects controllers.
		webRoutes.GET("/deletions/show/:id", web.DeletionRequestShow)
		webRoutes.GET("/deletions/review/:id", web.DeletionRequestReview)
		webRoutes.POST("/deletions/approve/:id", web.DeletionRequestApprove)
		webRoutes.POST("/deletions/cancel/:id", web.DeletionRequestCancel)
		webRoutes.GET("/deletions/", web.DeletionRequestIndex)

		// Dashboard
		webRoutes.GET("/dashboard", web.DashboardShow)

		// GenericFiles
		webRoutes.GET("/files", web.GenericFileIndex)
		webRoutes.GET("/files/show/:id", web.GenericFileShow)
		webRoutes.GET("/files/request_delete/:id", web.GenericFileRequestDelete)
		webRoutes.GET("/files/request_restore/:id", web.GenericFileRequestRestore)
		webRoutes.POST("/files/init_delete/:id", web.GenericFileInitDelete)
		webRoutes.POST("/files/init_restore/:id", web.GenericFileInitRestore)

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

		// IntellectualObjects
		webRoutes.GET("/objects", web.IntellectualObjectIndex)
		webRoutes.GET("/objects/show/:id", web.IntellectualObjectShow)
		webRoutes.GET("/objects/request_delete/:id", web.IntellectualObjectRequestDelete)
		webRoutes.POST("/objects/init_delete/:id", web.IntellectualObjectInitDelete)
		webRoutes.GET("/objects/request_restore/:id", web.IntellectualObjectRequestRestore)
		webRoutes.POST("/objects/init_restore/:id", web.IntellectualObjectInitRestore)
		webRoutes.GET("/objects/events/:id", web.IntellectualObjectEvents)
		webRoutes.GET("/objects/files/:id", web.IntellectualObjectFiles)

		// PremisEvents
		webRoutes.GET("/events", web.PremisEventIndex)
		webRoutes.GET("/events/show/:id", web.PremisEventShow)
		webRoutes.GET("/events/show_xhr/:id", web.PremisEventShowXHR)

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
		webRoutes.GET("/users/my_account", web.UserMyAccount)
		webRoutes.GET("/users/change_password/:id", web.UserShowChangePassword)
		webRoutes.POST("/users/change_password/:id", web.UserChangePassword)
		webRoutes.GET("/users/init_password_reset/:id", web.UserInitPasswordReset)
		webRoutes.GET("/users/complete_password_reset/:id", web.UserCompletePasswordReset)
		webRoutes.POST("/users/get_api_key/:id", web.UserGetAPIKey)

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
