package pgmodels_test

import (
	// "encoding/json"
	// "fmt"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositStats(t *testing.T) {
	db.LoadFixtures()

	date2010, err := time.Parse(time.RFC3339, "2010-01-01T00:00:00Z")
	require.Nil(t, err)
	date2030, err := time.Parse(time.RFC3339, "2030-01-01T00:00:00Z")
	require.Nil(t, err)

	// Nothing prior to 2010
	stats, err := pgmodels.DepositStatsSelect(3, constants.StorageOptionStandard, date2010)
	require.Nil(t, err)
	require.NotNil(t, stats)

	require.Equal(t, 1, len(stats))
	assert.Equal(t, "Total", stats[0].StorageOption)
	assert.EqualValues(t, 0, stats[0].ObjectCount)
	assert.EqualValues(t, 0, stats[0].FileCount)
	assert.EqualValues(t, 0, stats[0].TotalBytes)
	assert.EqualValues(t, 0, stats[0].MonthlyCost)

	// Prior to 2030
	// One storage option for one institution
	stats, err = pgmodels.DepositStatsSelect(3, constants.StorageOptionStandard, date2030)
	require.Nil(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, 4, len(stats))
	assert.Equal(t, "Standard", stats[0].StorageOption)
	assert.EqualValues(t, 30, stats[0].FileCount)
	assert.Equal(t, "Total", stats[1].StorageOption)
	assert.EqualValues(t, 30, stats[1].FileCount)

	// Prior to 2030
	// All storage options for one institution
	stats, err = pgmodels.DepositStatsSelect(3, "", date2030)
	require.Nil(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, 8, len(stats))
	assert.Equal(t, "Glacier-Deep-VA", stats[0].StorageOption)
	assert.EqualValues(t, 2, stats[0].FileCount)
	assert.Equal(t, "Standard", stats[1].StorageOption)
	assert.EqualValues(t, 30, stats[1].FileCount)

	// Prior to 2030
	// All storage options for all institutions
	stats, err = pgmodels.DepositStatsSelect(0, "", date2030)
	require.Nil(t, err)
	require.NotNil(t, stats)

	//js, _ := json.Marshal(stats)
	//fmt.Println(string(js))

	assert.Equal(t, 18, len(stats))
	assert.Equal(t, "Glacier-Deep-OH", stats[0].StorageOption)
	assert.EqualValues(t, 3, stats[0].FileCount)
	assert.Equal(t, "Glacier-OR", stats[1].StorageOption)
	assert.EqualValues(t, 2, stats[1].FileCount)
}
