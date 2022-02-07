package app

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/helpers"
	"github.com/APTrust/registry/middleware"
	admin_api "github.com/APTrust/registry/web/api/admin"
	common_api "github.com/APTrust/registry/web/api/common"
	"github.com/APTrust/registry/web/webui"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

// Run runs the Registry application. This is called from main() to start
// the app.
func Run() {
	r := InitAppEngine(false)
	ctx := common.Context()
	if ctx.Config.Cookies.HTTPSOnly && ctx.Config.Cookies.Domain != "localhost" {
		err := http.Serve(autocert.NewListener(ctx.Config.Cookies.Domain), r)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	} else {
		r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}
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
		"dateISO":        helpers.DateISO,
		"dateTimeISO":    helpers.DateTimeISO,
		"dateUS":         helpers.DateUS,
		"defaultString":  helpers.DefaultString,
		"dict":           helpers.Dict,
		"escapeAttr":     helpers.EscapeAttr,
		"escapeHTML":     helpers.EscapeHTML,
		"formatFloat":    helpers.FormatFloat,
		"humanSize":      helpers.HumanSize,
		"iconFor":        helpers.IconFor,
		"replace":        strings.Replace,
		"roleName":       helpers.RoleName,
		"strEq":          helpers.StrEq,
		"titleCase":      strings.Title,
		"toJSON":         helpers.ToJSON,
		"truncate":       helpers.Truncate,
		"truncateMiddle": helpers.TruncateMiddle,
		"truncateStart":  helpers.TruncateStart,
		"userCan":        helpers.UserCan,
		"yesNo":          helpers.YesNo,
	})

	// Load the view templates
	// If we're running from main, templates will come
	// from ./views. When running tests, templates come
	// from ../../views because http tests run from web
	// from ../../../views for member api and admin api
	// sub directory.
	if common.FileExists("./views") {
		router.LoadHTMLGlob("./views/**/*.html")
	} else if common.FileExists("../../views") {
		router.LoadHTMLGlob("../../views/**/*.html")
	} else {
		router.LoadHTMLGlob("../../../views/**/*.html")
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
	router.Use(middleware.CSRF())
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
	router.RedirectFixedPath = false

	router.Static("/static", "./static")
	router.Static("/favicon.ico", "./static/img/favicon.png")

	webRoutes := router.Group("/")
	{
		// Alerts
		webRoutes.GET("/alerts", webui.AlertIndex)
		webRoutes.GET("/alerts/show/:id/:user_id", webui.AlertShow)

		// Deletion Requests
		// Note that these routes are for read-only views.
		// Routes for initiating, approving and rejecting deletions
		// are in the GenericFiles and IntellectualObjects controllers.
		webRoutes.GET("/deletions/show/:id", webui.DeletionRequestShow)
		webRoutes.GET("/deletions/review/:id", webui.DeletionRequestReview)
		webRoutes.POST("/deletions/approve/:id", webui.DeletionRequestApprove)
		webRoutes.POST("/deletions/cancel/:id", webui.DeletionRequestCancel)
		webRoutes.GET("/deletions/", webui.DeletionRequestIndex)

		// Dashboard
		webRoutes.GET("/dashboard", webui.DashboardShow)

		// Deposit Report
		webRoutes.GET("/reports/deposits", webui.DepositReportShow)

		// GenericFiles
		webRoutes.GET("/files", webui.GenericFileIndex)
		webRoutes.GET("/files/show/:id", webui.GenericFileShow)
		webRoutes.GET("/files/request_delete/:id", webui.GenericFileRequestDelete)
		webRoutes.GET("/files/request_restore/:id", webui.GenericFileRequestRestore)
		webRoutes.POST("/files/init_delete/:id", webui.GenericFileInitDelete)
		webRoutes.POST("/files/init_restore/:id", webui.GenericFileInitRestore)

		// Institutions
		webRoutes.POST("/institutions/new", webui.InstitutionCreate)
		webRoutes.DELETE("/institutions/delete/:id", webui.InstitutionDelete)
		webRoutes.GET("/institutions/delete/:id", webui.InstitutionDelete)
		webRoutes.GET("/institutions/undelete/:id", webui.InstitutionUndelete)
		webRoutes.GET("/institutions", webui.InstitutionIndex)
		webRoutes.GET("/institutions/new", webui.InstitutionNew)
		webRoutes.GET("/institutions/show/:id", webui.InstitutionShow)
		webRoutes.GET("/institutions/edit/:id", webui.InstitutionEdit)
		webRoutes.PUT("/institutions/edit/:id", webui.InstitutionUpdate)
		webRoutes.POST("/institutions/edit/:id", webui.InstitutionUpdate)

		// IntellectualObjects
		webRoutes.GET("/objects", webui.IntellectualObjectIndex)
		webRoutes.GET("/objects/show/:id", webui.IntellectualObjectShow)
		webRoutes.GET("/objects/request_delete/:id", webui.IntellectualObjectRequestDelete)
		webRoutes.POST("/objects/init_delete/:id", webui.IntellectualObjectInitDelete)
		webRoutes.GET("/objects/request_restore/:id", webui.IntellectualObjectRequestRestore)
		webRoutes.POST("/objects/init_restore/:id", webui.IntellectualObjectInitRestore)
		webRoutes.GET("/objects/events/:id", webui.IntellectualObjectEvents)
		webRoutes.GET("/objects/files/:id", webui.IntellectualObjectFiles)

		// PremisEvents
		webRoutes.GET("/events", webui.PremisEventIndex)
		webRoutes.GET("/events/show/:id", webui.PremisEventShow)
		webRoutes.GET("/events/show_xhr/:id", webui.PremisEventShowXHR)

		// WorkItems - Web UI allows only list, show, and limited editing for admin only
		webRoutes.GET("/work_items", webui.WorkItemIndex)
		webRoutes.GET("/work_items/show/:id", webui.WorkItemShow)
		webRoutes.GET("/work_items/edit/:id", webui.WorkItemEdit)
		webRoutes.PUT("/work_items/edit/:id", webui.WorkItemUpdate)
		webRoutes.POST("/work_items/edit/:id", webui.WorkItemUpdate)
		webRoutes.PUT("/work_items/requeue/:id", webui.WorkItemRequeue)
		webRoutes.POST("/work_items/requeue/:id", webui.WorkItemRequeue)

		// Users
		webRoutes.POST("/users/new", webui.UserCreate)
		webRoutes.DELETE("/users/delete/:id", webui.UserDelete)
		webRoutes.POST("/users/delete/:id", webui.UserDelete)
		webRoutes.POST("/users/undelete/:id", webui.UserUndelete)
		webRoutes.PUT("/users/undelete/:id", webui.UserUndelete)
		webRoutes.GET("/users", webui.UserIndex)
		webRoutes.GET("/users/new", webui.UserNew)
		webRoutes.GET("/users/show/:id", webui.UserShow)
		webRoutes.GET("/users/edit/:id", webui.UserEdit)
		webRoutes.PUT("/users/edit/:id", webui.UserUpdate)
		webRoutes.POST("/users/edit/:id", webui.UserUpdate)
		webRoutes.GET("/users/my_account", webui.UserMyAccount)
		webRoutes.GET("/users/change_password/:id", webui.UserShowChangePassword)
		webRoutes.POST("/users/change_password/:id", webui.UserChangePassword)
		webRoutes.GET("/users/init_password_reset/:id", webui.UserInitPasswordReset)
		webRoutes.GET("/users/complete_password_reset/:id", webui.UserCompletePasswordReset)
		webRoutes.POST("/users/get_api_key/:id", webui.UserGetAPIKey)

		// User two-factor setup
		webRoutes.GET("/users/2fa_setup", webui.UserInit2FASetup)
		webRoutes.POST("/users/2fa_setup", webui.UserComplete2FASetup)
		webRoutes.POST("/users/confirm_phone", webui.UserConfirmPhone)
		webRoutes.POST("/users/backup_codes", webui.UserGenerateBackupCodes)

		// User two-factor login
		webRoutes.GET("/users/2fa_backup", webui.UserTwoFactorBackup)
		webRoutes.GET("/users/2fa_choose", webui.UserTwoFactorChoose)
		webRoutes.POST("/users/2fa_sms", webui.UserTwoFactorGenerateSMS)
		webRoutes.POST("/users/2fa_push", webui.UserTwoFactorPush)
		webRoutes.POST("/users/2fa_verify", webui.UserTwoFactorVerify)

		// User Sign In
		webRoutes.GET("/users/sign_in", webui.UserSignInShow)
		webRoutes.POST("/users/sign_in", webui.UserSignIn)
		webRoutes.GET("/users/sign_out", webui.UserSignOut) // should be delete?

		// Error page
		webRoutes.GET("/error", webui.ErrorShow)

	}

	// Root goes to sign-in page, which is a web route,
	// not an API route.
	router.GET("/", webui.UserSignInShow)

	// Member API routes. Note that the show routes for
	// GenericFiles, Institutions, IntellectualObjects and
	// PremisEvents end with *id instead of :id. This tells
	// julienschmidt/httprouter to cram everything after the
	// slash into the "id" parameter, which allows us to serve
	// files, institutions, objects and events by id or identifier.
	//
	// For example, assuming file with id 99 has identifier
	// "school.edu/bag_name/data/image.jpg", the following routes
	// return the same thing:
	//
	// /member-api/v3/files/99
	// /member-api/v3/files/school.edu/bag_name/data/image.jpg
	// /member-api/v3/files/school.edu%2Fbag_name%2Fdata%2Fimage.jpg
	//
	// If file identifier contains a question mark, it MUST be
	// url-encoded, or the router will interpret the ? as the
	// beginning of the query string. To be safe, we should always
	// url-encode the identifier. Many of them contain backticks,
	// quotes, parentheses, spaces, and all kinds of other garbage.
	//
	// Routes start with /member-api/v3
	memberAPI := router.Group(fmt.Sprintf("%sv3", constants.APIPrefixMember))
	{
		// Alerts
		// TODO: Delete this? Is there even a use case?
		memberAPI.GET("/alerts", common_api.AlertIndex)
		memberAPI.GET("/alerts/show/:id/:user_id", common_api.AlertShow)

		// Checksums
		memberAPI.GET("/checksums", common_api.ChecksumIndex)
		memberAPI.GET("/checksums/show/:id", common_api.ChecksumShow)

		// Deletion Requests
		// TODO: Should we really expose this through the API?
		memberAPI.GET("/deletions/show/:id", common_api.DeletionRequestShow)
		memberAPI.GET("/deletions", common_api.DeletionRequestIndex)

		// Generic Files
		memberAPI.GET("/files/show/*id", common_api.GenericFileShow)
		memberAPI.GET("/files", common_api.GenericFileIndex)

		// Intellectual Objects
		memberAPI.GET("/objects/show/*id", common_api.IntellectualObjectShow)
		memberAPI.GET("/objects", common_api.IntellectualObjectIndex)

		// Premis Events
		memberAPI.GET("/events/show/*id", common_api.PremisEventShow)
		memberAPI.GET("/events", common_api.PremisEventIndex)

		// Work Items
		memberAPI.GET("/items/show/:id", common_api.WorkItemShow)
		memberAPI.GET("/items", common_api.WorkItemIndex)

	}

	// Admin API is used by preservation-services.
	// Note that this group uses the same handlers
	// as the member API for some show and index routes.
	//
	// Routes start with /admin-api/v3
	adminAPI := router.Group(fmt.Sprintf("%sv3", constants.APIPrefixAdmin))
	{
		// Alerts
		// TODO: Delete this? Admin API doesn't really need it.
		adminAPI.GET("/alerts", common_api.AlertIndex)
		adminAPI.GET("/alerts/show/:id/:user_id", common_api.AlertShow)

		// Checksums
		adminAPI.GET("/checksums", common_api.ChecksumIndex)
		adminAPI.GET("/checksums/show/:id", common_api.ChecksumShow)

		// Deletion Requests
		// TODO: Does Admin API really need this?
		adminAPI.GET("/deletions/show/:id", common_api.DeletionRequestShow)
		adminAPI.GET("/deletions", common_api.DeletionRequestIndex)

		// Generic Files
		adminAPI.GET("/files/show/*id", common_api.GenericFileShow)
		adminAPI.GET("/files", common_api.GenericFileIndex)
		adminAPI.DELETE("/files/delete/:id", admin_api.GenericFileDelete)
		adminAPI.POST("/files/create/:institution_id", admin_api.GenericFileCreate)
		adminAPI.POST("/files/create_batch/:institution_id", admin_api.GenericFileCreateBatch)
		adminAPI.PUT("/files/update/:id", admin_api.GenericFileUpdate)

		// Institutions
		adminAPI.GET("/institutions", admin_api.InstitutionIndex)
		adminAPI.GET("/institutions/show/:id", admin_api.InstitutionShow)

		// Intellectual Objects
		adminAPI.GET("/objects/show/*id", common_api.IntellectualObjectShow)
		adminAPI.GET("/objects", common_api.IntellectualObjectIndex)
		adminAPI.POST("/objects/create/:institution_id", admin_api.IntellectualObjectCreate)
		adminAPI.PUT("/objects/update/:id", admin_api.IntellectualObjectUpdate)
		adminAPI.DELETE("/objects/delete/:id", admin_api.IntellectualObjectDelete)

		// Premis Events
		adminAPI.POST("/events/create", admin_api.PremisEventCreate)
		adminAPI.GET("/events/show/*id", common_api.PremisEventShow)
		adminAPI.GET("/events", common_api.PremisEventIndex)

		// Work Items
		// TODO: Add Create & Update
		adminAPI.POST("/items/create/:institution_id", admin_api.WorkItemCreate)
		adminAPI.PUT("/items/update/:id", admin_api.WorkItemUpdate)
		adminAPI.GET("/items/show/:id", common_api.WorkItemShow)
		adminAPI.GET("/items", common_api.WorkItemIndex)

	}
}
