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
