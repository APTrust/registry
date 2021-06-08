package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expected = []*pgmodels.ObjectStats{
	{1, 2233889, "application/binary"},
	{1, 78333, "audio/mp3"},
	{1, 988226, "text/plain"},
	{1, 44344, "text/sgml"},
	{1, 399288, "video/mp4"},
	{5, 3744080, ""},
}

func TestGetObjectStats(t *testing.T) {
	db.LoadFixtures()
	stats, err := pgmodels.GetObjectStats(3)
	require.Nil(t, err)
	assert.Equal(t, expected, stats)
}
