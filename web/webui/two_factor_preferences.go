package webui

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type TwoFactorPreferences struct {
	OldPhone  string
	NewPhone  string
	OldMethod string
	NewMethod string
	User      *pgmodels.User
}

func NewTwoFactorPreferences(req *Request) (*TwoFactorPreferences, error) {
	oldPhone := req.CurrentUser.PhoneNumber
	oldMethod := req.CurrentUser.AuthyStatus

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
		OldPhone:  oldPhone,
		NewPhone:  user.PhoneNumber,
		OldMethod: oldMethod,
		NewMethod: user.AuthyStatus,
		User:      user,
	}

	return prefs, nil
}

func (p *TwoFactorPreferences) PhoneChanged() bool {
	return p.OldPhone != p.NewPhone
}

func (p *TwoFactorPreferences) MethodChanged() bool {
	return p.OldMethod != p.NewMethod
}

func (p *TwoFactorPreferences) NeedsConfirmation() bool {
	return p.PhoneChanged() || p.MethodChanged()
}

func (p *TwoFactorPreferences) NothingChanged() bool {
	return !p.PhoneChanged() && !p.MethodChanged()
}

func (p *TwoFactorPreferences) DoNotUseTwoFactor() bool {
	return p.NewMethod == constants.TwoFactorNone
}

func (p *TwoFactorPreferences) UsePasskey() bool {
	return p.NewMethod == constants.TwoFactorPasskey
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

func (p *TwoFactorPreferences) NeedsSMSConfirmation() bool {
	return p.NeedsConfirmation() && p.NewMethod == constants.TwoFactorSMS
}
