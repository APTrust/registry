package webui

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type TwoFactorPreferences struct {
	OldPhone         string
	NewPhone         string
	OldMethod        string
	NewMethod        string
	OldAuthAppMethod string
	NewAuthAppMethod string
	User             *pgmodels.User
}

func NewTwoFactorPreferences(req *Request) (*TwoFactorPreferences, error) {
	oldPhone := req.CurrentUser.PhoneNumber
	oldMethod := req.CurrentUser.AuthyStatus
	oldAuthAppMethod := constants.TwoFactorNone // req.CurrentUser.UseAuthenticatorApp

	// Get phone and authy data submitted in the form.
	user := &pgmodels.User{}
	err := req.GinContext.ShouldBind(user)
	if err != nil {
		return nil, err
	}

	// Make sure phone is formatted the way Authy & SMS/SNS like it,
	// with leading + and country code. https://trello.com/c/QLMjQiyj
	user.ReformatPhone()

	prefs := &TwoFactorPreferences{
		OldPhone:         oldPhone,
		NewPhone:         user.PhoneNumber,
		OldMethod:        oldMethod,
		NewMethod:        user.AuthyStatus,
		User:             user,
		OldAuthAppMethod: oldAuthAppMethod,
		NewAuthAppMethod: constants.TwoFactorNone,
	}

	return prefs, nil
}

func (p *TwoFactorPreferences) PhoneChanged() bool {
	return p.OldPhone != p.NewPhone
}

func (p *TwoFactorPreferences) MethodChanged() bool {
	return p.OldMethod != p.NewMethod
}

func (p *TwoFactorPreferences) AuthAppMethodChanged() bool {
	return p.OldAuthAppMethod != p.NewAuthAppMethod
}

func (p *TwoFactorPreferences) NeedsConfirmation() bool {
	return p.PhoneChanged() || p.MethodChanged() || p.AuthAppMethodChanged()
}

func (p *TwoFactorPreferences) NothingChanged() bool {
	return !p.PhoneChanged() && !p.MethodChanged() && !p.AuthAppMethodChanged()
}

func (p *TwoFactorPreferences) DoNotUseTwoFactor() bool {
	return p.NewMethod == constants.TwoFactorNone
}

func (p *TwoFactorPreferences) UseAuthenticatorApp() bool {
	return p.NewMethod == constants.TwoFactorTOTP
}

func (p *TwoFactorPreferences) UseAuthy() bool {
	return p.NewMethod == constants.TwoFactorAuthy
}

func (p *TwoFactorPreferences) UseSMS() bool {
	return p.NewMethod == constants.TwoFactorSMS
}

func (p *TwoFactorPreferences) NeedsAuthyRegistration() bool {
	return p.NewMethod == constants.TwoFactorAuthy && p.User.AuthyID == ""
}

func (p *TwoFactorPreferences) NeedsAuthyConfirmation() bool {
	return p.NeedsConfirmation() && p.NewMethod == constants.TwoFactorAuthy
}

func (p *TwoFactorPreferences) NeedsAuthenticatorAppRegistration() bool {
	return true
	// return p.NewAuthAppMethod == constants.TwoFactorTOTP && p.User.TOTPSecret == ""
}

func (p *TwoFactorPreferences) NeedsAuthenticatorAppConfirmation() bool {
	return p.NeedsConfirmation() && p.NewMethod == constants.TwoFactorTOTP
}

func (p *TwoFactorPreferences) NeedsSMSConfirmation() bool {
	return p.NeedsConfirmation() && p.NewMethod == constants.TwoFactorSMS
}
