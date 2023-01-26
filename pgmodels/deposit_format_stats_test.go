package pgmodels_test

import (
	"testing"

	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositFormatStats(t *testing.T) {
	db.LoadFixtures()

	// By InstID + Obj ID
	stats, err := pgmodels.DepositFormatStatsSelect(3, 6)
	require.Nil(t, err)
	require.NotNil(t, stats)

	//j, _ := json.MarshalIndent(stats, "", "  ")
	//fmt.Println(string(j))

	require.Equal(t, 5, len(stats))
	assert.EqualValues(t, 24, pgmodels.StatsByFormat(stats, "Total").FileCount)

	// By Obj ID, without Inst ID
	stats, err = pgmodels.DepositFormatStatsSelect(0, 6)
	require.Nil(t, err)
	require.NotNil(t, stats)

	require.Equal(t, 5, len(stats))
	assert.EqualValues(t, 24, pgmodels.StatsByFormat(stats, "Total").FileCount)

	// By Inst ID only - all files for inst 3
	stats, err = pgmodels.DepositFormatStatsSelect(3, 0)
	require.Nil(t, err)
	require.NotNil(t, stats)

	require.Equal(t, 9, len(stats))
	assert.EqualValues(t, 35, pgmodels.StatsByFormat(stats, "Total").FileCount)
}
