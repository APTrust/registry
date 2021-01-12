package common_test

import (
	"path"
	"strings"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectRoot(t *testing.T) {
	assert.NotEmpty(t, common.ProjectRoot)
}

func TestLoadRelativeFile(t *testing.T) {
	file := path.Join("common", "config.go")
	data, err := common.LoadRelativeFile(file)
	require.Nil(t, err)
	assert.NotEmpty(t, data)
	code := string(data)
	assert.True(t, strings.HasPrefix(code, "package common"))
}

func TestExpandTilde(t *testing.T) {
	expanded, err := common.ExpandTilde("~/tmp")
	assert.Nil(t, err)
	assert.True(t, len(expanded) > 6)
	assert.True(t, strings.HasSuffix(expanded, "tmp"))

	expanded, err = common.ExpandTilde("/nothing/to/expand")
	assert.Nil(t, err)
	assert.Equal(t, "/nothing/to/expand", expanded)
}

func TestHash(t *testing.T) {
	p := "be4296045a168d1c2a484b625bf67477dcf748c86de3f2a94a14ba9eefb53669"
	assert.Equal(t, p, common.Hash("password"))
}

func TestAESEncryption(t *testing.T) {
	config := common.NewConfig()

	// Opening words of one of the best noirs ever written
	plaintext := "The Academic Preservation Trust (APTrust) is committed to the creation and management of a sustainable environment for digital preservation."
	expected := "0000000000000000000000001c628c322eb59be9dc62e19707314b8febc459b766306d2d7158018680ce728d9754ffb462da4e8ef921f968a3656f10e4fc3259ee6bb620481fdec399a61412be2824c709233e8b7e934ea2b1b189a97c8cabcbd993a157ea9f5a7d94c40f8947d0795a380389a0321a78c7d05027f98719ca8f7c6917049607de2c63daded593c0bee2f7aca7833942201218a746aa06b70b9dc43205d2aad1e65f"

	hexCipher, err := common.EncryptAES(config.AESKey, plaintext)
	require.Nil(t, err)
	assert.Equal(t, expected, hexCipher)

	decrypted, err := common.DecryptAES(config.AESKey, hexCipher)
	require.Nil(t, err)
	assert.Equal(t, plaintext, decrypted)
}
