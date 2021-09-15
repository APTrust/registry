package webui

import (
	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

type TwoFactorPreferences struct {
	OldPhone  string
	NewPhone  string
	OldMethod string
	NewMethod string
	user      *pgmodels.User
}

func NewTwoFactorPreferences(req *Request) (*TwoFactorPreferences, error) {
	user := req.CurrentUser
	oldPhone := user.PhoneNumber
	oldMethod := user.AuthyStatus
	err := req.GinContext.ShouldBind(user)
	if err != nil {
		return nil, err
	}
	valError := user.Validate()
	if valError != nil && len(valError.Errors) > 0 {
		common.Context().Log.Error().Msgf("User validation error while completing two-factor setup: %s", valError.Error())
		return nil, valError
	}

	prefs := &TwoFactorPreferences{
		OldPhone:  oldPhone,
		NewPhone:  user.PhoneNumber,
		OldMethod: oldMethod,
		NewMethod: user.AuthyStatus,
		user:      user,
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

func (p *TwoFactorPreferences) UseAuthy() bool {
	return p.NewMethod == constants.TwoFactorAuthy
}

func (p *TwoFactorPreferences) UseSMS() bool {
	return p.NewMethod == constants.TwoFactorSMS
}

func (p *TwoFactorPreferences) NeedsAuthyRegistration() bool {
	return p.NewMethod == constants.TwoFactorAuthy && p.user.AuthyID == ""
}

func (p *TwoFactorPreferences) NeedsAuthyConfirmation() bool {
	return p.NeedsConfirmation() && p.NewMethod == constants.TwoFactorAuthy
}

func (p *TwoFactorPreferences) NeedsSMSConfirmation() bool {
	return p.NeedsConfirmation() && p.NewMethod == constants.TwoFactorSMS
}
