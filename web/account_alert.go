package web

import (
	"bytes"
	"fmt"
	"time"

	"github.com/APTrust/registry/common"
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
	return createAlert(alert, templateName, alertData)
}

// CreatePasswordResetAlert creates an alert telling a user to
// reset their password. The alert's email body includes a link
// with a password-reset token. These alerts should be generated
// when a user clicks the "Forgot Password" link on the login page
// and when an admin creates a new user account. In each case, the
// user clicks the email link and then has to set a new password.
func CreatePasswordResetAlert() {

}

// createAlert adds customized text to the alert and saves it in the
// database. Param templateName is the name of the text template used
// to construct the alert message. Param alertData is the custom data
// to put into the template.
//
// This returns the alert with a non-zero ID (since it saves it) and
// an error if there's a problem with the template or the save.
func createAlert(alert *pgmodels.Alert, templateName string, alertData map[string]interface{}) (*pgmodels.Alert, error) {

	// Create the alert text from the template...
	tmpl := common.TextTemplates[templateName]
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, alertData)
	if err != nil {
		return nil, err
	}

	// Set the alert text and save...
	alert.Content = buf.String()
	err = alert.Save()
	if err != nil {
		return nil, err
	}

	// Show the alert text in dev and test consoles,
	// so we don't have to look it up in the DB.
	// For dev/test, we need to see the review and
	// confirmation URLS in this alert so we can
	// review and test them.
	envName := common.Context().Config.EnvName
	if envName == "dev" || envName == "test" {
		fmt.Println("***********************")
		fmt.Println(alert.Content)
		fmt.Println("***********************")
	}

	return alert, err
}
