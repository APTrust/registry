package webui_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/web/testutil"
	"github.com/APTrust/registry/web/webui"
	"github.com/stretchr/testify/assert"
)

func TestTwoFactorPreferences(t *testing.T) {
	prefs := &webui.TwoFactorPreferences{
		OldPhone:  "",
		NewPhone:  "",
		OldMethod: "",
		NewMethod: "",
		User:      testutil.Inst1User,
	}

	assert.False(t, prefs.PhoneChanged())
	prefs.NewPhone = "12223334444"
	assert.True(t, prefs.PhoneChanged())

	assert.False(t, prefs.MethodChanged())
	prefs.NewMethod = constants.TwoFactorAuthy
	assert.True(t, prefs.MethodChanged())

	prefs.NewMethod = constants.TwoFactorNone
	assert.True(t, prefs.DoNotUseTwoFactor())
	assert.False(t, prefs.UseAuthy())
	assert.False(t, prefs.UseSMS())

	prefs.NewMethod = constants.TwoFactorSMS
	assert.False(t, prefs.DoNotUseTwoFactor())
	assert.False(t, prefs.UseAuthy())
	assert.True(t, prefs.UseSMS())

	prefs.NewMethod = constants.TwoFactorAuthy
	assert.False(t, prefs.DoNotUseTwoFactor())
	assert.True(t, prefs.UseAuthy())
	assert.False(t, prefs.UseSMS())

	prefs.NewMethod = constants.TwoFactorAuthy
	prefs.User.AuthyID = ""
	assert.True(t, prefs.NeedsAuthyRegistration())

	prefs.NewMethod = constants.TwoFactorSMS
	prefs.User.AuthyID = ""
	assert.False(t, prefs.NeedsAuthyRegistration())

	prefs.NewMethod = constants.TwoFactorAuthy
	prefs.User.AuthyID = "12345"
	assert.False(t, prefs.NeedsAuthyRegistration())

	assert.True(t, prefs.NeedsAuthyConfirmation())
	prefs.NewMethod = constants.TwoFactorSMS
	assert.False(t, prefs.NeedsAuthyConfirmation())

	assert.True(t, prefs.NeedsSMSConfirmation())
}
