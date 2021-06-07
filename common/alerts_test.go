package common_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var templateNames = []string{
	"alerts/deletion_cancelled.txt",
	"alerts/deletion_completed.txt",
	"alerts/deletion_confirmed.txt",
	"alerts/deletion_requested.txt",
	"alerts/failed_fixity.txt",
	"alerts/restoration_completed.txt",
}

// Make sure these templates are loaded, and that they have
// content. Unlike our HTML templates, if we include "define"
// directives in the text templates, they will silently fail
// and return no content after execution.
func TestAlertTemplates(t *testing.T) {
	require.NotNil(t, common.TextTemplates)
	for _, name := range templateNames {
		assert.NotNil(t, common.TextTemplates[name])
	}

	name := "Spongebob"
	link := "https://example.com/confirm?token=ABCD"
	data := map[string]string{
		"requesterName":     name,
		"deletionReviewURL": link,
	}
	tmpl := common.TextTemplates["alerts/deletion_requested.txt"]
	require.NotNil(t, tmpl)
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	require.Nil(t, err)

	content := buf.String()
	assert.True(t, len(content) > 100)
	assert.True(t, strings.Contains(content, name))
	assert.True(t, strings.Contains(content, link))
}
