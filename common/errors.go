package common

import (
	"errors"
)

var ErrInvalidLogin = errors.New("invalid login")
var ErrAccountDeactivated = errors.New("account deactivated")
var ErrPermissionDenied = errors.New("permission denied")
var ErrNotSupported = errors.New("operation not supported")
