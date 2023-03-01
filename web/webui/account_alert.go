package webui

import (
	"fmt"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// CreatePasswordChangedAlert creates an alert telling a user that
// their password has been changed. It should say who changed the
// password, from where (IP address) and when.
func CreatePasswordChangedAlert(req *Request, userToEdit *pgmodels.User) (*pgmodels.Alert, error) {
	templateName := "alerts/password_changed.txt"
	alertData := map[string]interface{}{
		"registryURL": req.BaseURL(),
		"userName":    req.CurrentUser.Name,
		"changeDate":  time.Now().Format(time.RFC3339),
		"userAgent":   req.GinContext.GetHeader("User-Agent"),
		"ipAddress":   req.GinContext.ClientIP(),
	}
	recipients := []*pgmodels.User{userToEdit}
	alert := &pgmodels.Alert{
		InstitutionID: userToEdit.InstitutionID,
		Type:          constants.AlertPasswordChanged,
		Subject:       "Your APTrust password has been changed",
		CreatedAt:     time.Now().UTC(),
		Users:         recipients,
	}
	return pgmodels.CreateAlert(alert, templateName, alertData)
}

// CreatePasswordResetAlert creates an alert telling a user to
// reset their password. The alert's email body includes a link
// with a password-reset token. These alerts should be generated
// when a user clicks the "Forgot Password" link on the login page
// and when an admin creates a new user account. In each case, the
// user clicks the email link and then has to set a new password.
//
// Note that the token in the URL is random but unencrypted. The one
// in the DB is encrypted, so stealing it will do an attacker no good.
func CreatePasswordResetAlert(req *Request, userToEdit *pgmodels.User, token string) (*pgmodels.Alert, error) {
	templateName := "alerts/reset_password.txt"
	alertData := map[string]interface{}{
		"passwordResetURL":   fmt.Sprintf("%s/users/complete_password_reset/%d", req.BaseURL(), userToEdit.ID),
		"passwordResetToken": token,
	}
	recipients := []*pgmodels.User{userToEdit}
	alert := &pgmodels.Alert{
		InstitutionID: userToEdit.InstitutionID,
		Type:          constants.AlertPasswordReset,
		Subject:       "Reset your APTrust password",
		CreatedAt:     time.Now().UTC(),
		Users:         recipients,
	}
	alert, err := pgmodels.CreateAlert(alert, templateName, alertData)

	if err == nil {
		userToEdit.ResetPasswordSentAt = time.Now().UTC()
		err = userToEdit.Save()
	}

	return alert, err
}

// CreateNewAccountAlert creates an alert telling a newly added user that
// they now have an APTrust account. The alert includes a link for
// the user to set a custom password.
func CreateNewAccountAlert(req *Request, newUser *pgmodels.User, token string) (*pgmodels.Alert, error) {
	templateName := "alerts/welcome.txt"
	alertData := map[string]interface{}{
		"passwordResetURL":   fmt.Sprintf("%s/users/complete_password_reset/%d", req.BaseURL(), newUser.ID),
		"passwordResetToken": token,
		"adminName":          req.CurrentUser.Name,
	}
	recipients := []*pgmodels.User{newUser}
	alert := &pgmodels.Alert{
		InstitutionID: newUser.InstitutionID,
		Type:          constants.AlertWelcome,
		Subject:       "Your New APTrust Registry Account",
		CreatedAt:     time.Now().UTC(),
		Users:         recipients,
	}
	alert, err := pgmodels.CreateAlert(alert, templateName, alertData)

	if err == nil {
		newUser.ResetPasswordSentAt = time.Now().UTC()
		err = newUser.Save()
	}

	return alert, err
}
