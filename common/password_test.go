package common_test

import (
	"strings"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Sample encrypted password:
//
// $2a$10$Y6xQxU91zUE1WnThPxdGu.SbA/aruuheTxvBB6FO1N6LTM6ARXo5u
//
// Anatomy:
//
// $2a                                 - The bcrypt algorithm version
// $10                                 - The cost / number of stretches
// Y6xQxU91zUE1WnThPxdGu               - The salt
// .SbA/aruuheTxvBB6FO1N6LTM6ARXo5u    - Encrypted password w/salt & stretches

func TestEncryptPassword(t *testing.T) {
	enc, err := common.EncryptPassword("seekrit!")
	require.Nil(t, err)

	// Encrypted password always begins with algorithm version (2a)
	// and cost (10). The dollar signs are delimiters.
	assert.True(t, strings.HasPrefix(enc, "$2a$10$"))

	// Encrypted password has 29 byte prefix and 31 byte value.
	// Prefix is version + cost + salt, delimited by dollar signs.
	assert.Equal(t, 60, len(enc))
}

func TestComparePasswords(t *testing.T) {
	enc, err := common.EncryptPassword("homer")
	require.Nil(t, err)
	assert.True(t, common.ComparePasswords(enc, "homer"))

	enc, err = common.EncryptPassword("marge")
	require.Nil(t, err)
	assert.True(t, common.ComparePasswords(enc, "marge"))
}

// Note that in our fixture data under db/fixtures/users.csv, all test
// users have the password "password". This password was encrypted in
// Pharos under Devise, and we should be able to decrypt it with our
// code. Doing so means we can migrate to the new system without
// forcing all users to reset their passwords.
func TestLegacyPassword(t *testing.T) {
	plaintext := "password"
	hashed := "$2a$10$7aoot2KFFqikpTYVEbErYOxZijCHDPvqT4OMoFwdmsYBE9SK2PibC"
	assert.True(t, common.ComparePasswords(hashed, plaintext))
	assert.False(t, common.ComparePasswords(hashed, "should-not-match"))
}

func TestRandomToken(t *testing.T) {
	seen := make(map[string]bool)
	for i := 1; i < 5; i++ {
		token := common.RandomToken()
		assert.False(t, seen[token])
		seen[token] = true
	}
}

func TestLooksEncrypted(t *testing.T) {
	encrypted, err := common.EncryptPassword("Ned Flanders")
	require.Nil(t, err)
	assert.True(t, common.LooksEncrypted(encrypted))
	assert.False(t, common.LooksEncrypted("Barney Gumble"))
}

func TestPasswordMeetsRequirements(t *testing.T) {
	// Too short
	assert.False(t, common.PasswordMeetsRequirements("aBc1"))

	// No uppercase
	assert.False(t, common.PasswordMeetsRequirements("abc12345678"))

	// No lowercase
	assert.False(t, common.PasswordMeetsRequirements("ABC12345678"))

	// No numeric
	assert.False(t, common.PasswordMeetsRequirements("abcABCxyzXYZ"))

	// Goldilocks! Just right!
	assert.True(t, common.PasswordMeetsRequirements("abc123XYZ"))
	assert.True(t, common.PasswordMeetsRequirements("IAmOk110"))
}
