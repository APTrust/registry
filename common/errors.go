package common

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNotSignedIn means user has not signed in.
var ErrNotSignedIn = errors.New("user is not signed in")

// ErrInvalidLogin means the user supplied the wrong login name
// or password while trying to sign in.
var ErrInvalidLogin = errors.New("invalid login or password")

// ErrInvalidAPICredentials means the user supplied the wrong email
// or API token while trying to access a REST API endpoint.
var ErrInvalidAPICredentials = errors.New("invalid api credentials")

// ErrAccountDeactivated means the user logged in to a deactivated
// account.
var ErrAccountDeactivated = errors.New("account deactivated")

// ErrPermissionDenied means the user tried to access a resource
// withouth sufficient permission.
var ErrPermissionDenied = errors.New("permission denied")

// ErrNotSupported is an internal error that occurs, for example,
// when we try to delete an object that does not support deletion.
// This represents a programmer error and should not occur.
var ErrNotSupported = errors.New("operation not supported")

// ErrDecodeCookie occurs when we get a bad authentication cookie.
// We consider it a 400/Bad Request because we don't set bad cookies.
var ErrDecodeCookie = errors.New("error decoding cookie")

// ErrInvalidParam means the HTTP request contained an invalid
// parameter.
var ErrInvalidParam = errors.New("invalid parameter")

// ErrWrongDataType occurs when the user submits data of the wrong type,
// such string data that cannot be converted to a number, bool, date,
// or whatever type the application is expecting.
var ErrWrongDataType = errors.New("wrong data type")

// ErrParentRecordNotFound occurs when we cannot find the parent
// record required to check a user's permission. For example, when
// a user requests a Checksum, we first need to know if the user
// is allowed to access the Checksum's parent, which is a Generic
// File. If that record is missing, we get this error.
var ErrParentRecordNotFound = errors.New("parent record not found")

// ErrResourcePermission occurs when Authorization middleware cannot
// determine which permission is required to access the specified resoruce.
var ErrResourcePermission = errors.New("cannot determine permission type for requested requested resource")

// ErrPendingWorkItems occurs when a user wants to restore or delete an
// object or file but the WorkItems list shows other operations are pending
// on that item. For example, we can't delete or restore an object or file
// while another version of that object/file is pending ingest. Doing so
// would cause newly ingested files to be deleted as soon as they're sent
// to preservation, or would cause a restoration to contain a mix of new
// and old versions of a bag's files.
var ErrPendingWorkItems = errors.New("task cannot be completed because this object has pending work items")

// ErrInvalidToken means that the token presented for an action like
// password reset or deletion confirmation does not match the encrypted
// token in the database. When this error occurs, the user may not
// proceed with the requested action.
var ErrInvalidToken = errors.New("invalid token")

// ErrInvalidCSRFToken is specifically for POST/PUT/DELETE where we
// have a missing or invalid CSRF token.
var ErrInvalidCSRFToken = errors.New("invalid csrf token")

// ErrMissingReferer means the request had no http referer header,
// so we can't tell where it came from. We can't tell for sure if
// it's a cross-origin attack, so we throw an error. This happens
// only on unsafe methods (POST, PUT, DELETE). See the CSRF middleware.
var ErrMissingReferer = errors.New("csrf error: missing http referer")

// ErrMissingReferer means we got a cross-site request on an unsafe
// methods such as POST, PUT, DELETE. See the CSRF middleware.
var ErrCrossOriginReferer = errors.New("csrf error: cross origin request forbidden")

// ErrPasswordReqs indicates that the password the user is trying to set
// does not meet minimum requirements.
var ErrPasswordReqs = errors.New("password does not meet minimum requirements")

// ErrAlreadyHasAuthyID occurs when we try to register a user with Authy
// and they already have an Authy ID.
var ErrAlreadyHasAuthyID = errors.New("user is already registered with authy")

// ErrNoAuthyID occurs when a user requests two-factor login via Authy
// but does not have an Authy ID.
var ErrNoAuthyID = errors.New("user does not have an authy id")

// ErrWrongAPI occurs when a non-admin user tries to access the admin API.
// While the member and admin APIs share some common handlers, and members
// do technically have access to a number of read-only operations in both
// APIs, we don't want members to get in the habit of accessing the wrong
// endpoints.
var ErrWrongAPI = errors.New("non-admins must use the member api")

// ErrIDMismatch can indicate that the specified resource does not belong to
// the specified institition, or (more likely) that resource/institution IDs
// in a URL do not match resource/institution IDs submitted in the PUT/POST
// body of an HTML request.
var ErrIDMismatch = errors.New("resource or institution id mismatch")

// ErrInstIDChange occurs when someone tries to change the institution id
// of a resource.
var ErrInstIDChange = errors.New("institution id cannot change")

// ErrIdentifierChange occurs when someone tries to change the identifier
// of a resource.
var ErrIdentifierChange = errors.New("identifier cannot change")

// ErrStorageOptionChange indicates someone tried to change the storage
// option on an active object. (It's OK to change storage option on a
// deleted object that you are re-depositing.)
var ErrStorageOptionChange = errors.New("cannot change storage option on active object")

// ErrActiveFiles occurs when we try to delete an IntellectualObject
// that has active files. All of the object's files should first be
// deleted (set to State = "D").
var ErrActiveFiles = errors.New("cannot delete object with active files")

// ErrInternal is a runtime error that is not the user's fault, hence
// probably the programmer's fault.
var ErrInternal = errors.New("internal server error")

// ErrSubclassMustImplement indicates that a subclass has not implemented
// a method inherited from a base class.
var ErrSubclassMustImplement = errors.New("subclass must implement this method")

// ErrCountTypeNotSupported indicates that we cannot get counts for the
// given type from the *_counts views. Use a regular SQL count() instead.
var ErrCountTypeNotSupported = errors.New("type is not supported for view count")

type ValidationError struct {
	Errors map[string]string
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make(map[string]string),
	}
}

func (v *ValidationError) Error() string {
	if v.Errors == nil {
		return ""
	}
	errs := make([]string, len(v.Errors))
	i := 0
	for field, errMsg := range v.Errors {
		errs[i] = fmt.Sprintf("%s: %s", field, errMsg)
		i++
	}
	return strings.Join(errs, "\n")
}
