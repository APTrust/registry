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

	prefs.NewMethod = constants.TwoFactorNone
	assert.True(t, prefs.DoNotUseTwoFactor())
	assert.False(t, prefs.UseSMS())

	prefs.NewMethod = constants.TwoFactorSMS
	assert.False(t, prefs.DoNotUseTwoFactor())
	assert.True(t, prefs.UseSMS())

	assert.True(t, prefs.NeedsSMSConfirmation())
}
