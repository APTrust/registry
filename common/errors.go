package common

import (
	"errors"
)

// ErrNotSignedIn means user has not signed in.
var ErrNotSignedIn = errors.New("user is not signed in")

// ErrInvalidLogin means the user supplied the wrong login name
// or password while trying to sign in.
var ErrInvalidLogin = errors.New("invalid login or password")

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
