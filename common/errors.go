package common

import (
	"errors"
)

var ErrInvalidLogin = errors.New("invalid login or password")
var ErrAccountDeactivated = errors.New("account deactivated")
var ErrPermissionDenied = errors.New("permission denied")
var ErrNotSupported = errors.New("operation not supported")

var ErrDecodeCookie = errors.New("error decoding cookie")
var ErrWrongDataType = errors.New("wrong data type")
