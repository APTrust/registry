package web

import (
	"bytes"
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// TODO: Refactor Deletion.createAlert() so it's reusable.

// CreatePasswordChangedAlert creates an alert telling a user that
// their password has been changed. It should say who changed the
// password, from where (IP address) and when.
func CreatePasswordChangedAlert() {

}

// CreatePasswordResetAlert creates an alert telling a user to
// reset their password. The alert's email body includes a link
// with a password-reset token. These alerts should be generated
// when a user clicks the "Forgot Password" link on the login page
// and when an admin creates a new user account. In each case, the
// user clicks the email link and then has to set a new password.
func CreatePasswordResetAlert() {

}
