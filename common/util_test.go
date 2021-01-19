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
