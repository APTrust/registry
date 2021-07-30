package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expected = []*pgmodels.ObjectStats{
	{1, 11169445000, "application/binary"},
	{1, 391665000, "audio/mp3"},
	{1, 4941130000, "text/plain"},
	{1, 221720000, "text/sgml"},
	{1, 1996440000, "video/mp4"},
	{5, 18720400000, ""},
}

func TestGetObjectStats(t *testing.T) {
	db.LoadFixtures()
	stats, err := pgmodels.GetObjectStats(3)
	require.Nil(t, err)
	assert.Equal(t, expected, stats)
}
