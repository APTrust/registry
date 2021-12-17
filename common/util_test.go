package common_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
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

func TestFileExists(t *testing.T) {
	thisFile := path.Join(common.ProjectRoot(), "common", "util_test.go")
	assert.True(t, common.FileExists(thisFile))

	doesNotExist := path.Join(common.ProjectRoot(), "common", "nope.xyz")
	assert.False(t, common.FileExists(doesNotExist))
}

func TestCopyFile(t *testing.T) {
	src := path.Join(common.ProjectRoot(), "common", "util_test.go")
	dst := path.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().Unix()))
	defer os.Remove(dst)

	err := common.CopyFile(src, dst, 0644)
	require.Nil(t, err)

	finfo, err := os.Stat(dst)
	require.Nil(t, err)
	assert.Equal(t, os.FileMode(0644), finfo.Mode())
}

func TestListIsEmpty(t *testing.T) {
	empty1 := []string{}
	empty2 := []string{"", "", ""}
	notEmpty := []string{"", "yes", ""}

	assert.True(t, common.ListIsEmpty(empty1))
	assert.True(t, common.ListIsEmpty(empty2))
	assert.False(t, common.ListIsEmpty(notEmpty))
	assert.True(t, common.ListIsEmpty(nil))
}

func TestInterfaceList(t *testing.T) {
	list := common.InterfaceList(constants.States)
	require.Equal(t, 2, len(list))
	assert.Equal(t, "A", list[0].(string))
	assert.Equal(t, "D", list[1].(string))
}

func TestSplitCamelCase(t *testing.T) {
	assert.Equal(t, []string{"Institution", "Index"}, common.SplitCamelCase("InstitutionIndex", -1))
	assert.Equal(t, []string{"Currency", "USD"}, common.SplitCamelCase("CurrencyUSD", -1))

	// Split into all words
	assert.Equal(t, []string{"User", "Sign", "In"}, common.SplitCamelCase("UserSignIn", -1))

	// Split into 2 words max (first two)
	assert.Equal(t, []string{"User", "SignIn"}, common.SplitCamelCase("UserSignIn", 2))
}

func TestToHumanSize(t *testing.T) {
	assert.Equal(t, "389.8 kB", common.ToHumanSize(389778, 1000))
	assert.Equal(t, "380.6 kB", common.ToHumanSize(389778, 1024))
	assert.Equal(t, "3.9 GB", common.ToHumanSize(3897784432, 1000))
	assert.Equal(t, "3.6 GB", common.ToHumanSize(3897784432, 1024))
}

func TestCountryCodeAndPhone(t *testing.T) {
	cc, phone, err := common.CountryCodeAndPhone("+12125551212")
	require.Nil(t, err)
	assert.Equal(t, int32(1), cc)
	assert.Equal(t, "2125551212", phone)

	_, _, err = common.CountryCodeAndPhone("invalid number")
	assert.NotNil(t, err)
}

func TestIsEmptyString(t *testing.T) {
	assert.True(t, common.IsEmptyString(""))
	assert.True(t, common.IsEmptyString("   "))
	assert.True(t, common.IsEmptyString(" \t \n "))
	assert.False(t, common.IsEmptyString("not empty"))
	assert.False(t, common.IsEmptyString("  not empty  "))
}

func TestSanitizeIdentifier(t *testing.T) {
	assert.Equal(t, "this_34_xxyadda.yadda", common.SanitizeIdentifier("this_34_!;-xx{yadda.yadda}"))
	assert.Equal(t, `"users"."email"`, common.SanitizeIdentifier(`"users"."email"`))
}

func TestLooksLikeUUID(t *testing.T) {
	assert.True(t, common.LooksLikeUUID("1552abf5-28f3-46a5-ba63-95302d08e209"))
	assert.True(t, common.LooksLikeUUID("88198c5a-ec91-4ce1-bfcc-0f607ebdcca3"))
	assert.True(t, common.LooksLikeUUID("88198C5A-EC91-4CE1-BFCC-0F607EBDCCA3"))
	assert.False(t, common.LooksLikeUUID("88198c5a-ec91-4ce1-bfcc-0f607ebdccx3"))
	assert.False(t, common.LooksLikeUUID("88198c5a-ec91-4ce1-bfcc-0f6c"))
	assert.False(t, common.LooksLikeUUID(""))
}
