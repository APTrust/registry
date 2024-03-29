package pgmodels_test

import (
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertView(t *testing.T) {
	db.LoadFixtures()
	alertView, err := pgmodels.AlertViewForUser(1, 1)
	require.Nil(t, err)
	require.NotNil(t, alertView)

	assert.Equal(t, int64(1), alertView.GetID())
	assert.Equal(t, common.ErrNotSupported, alertView.Save())

	query := pgmodels.NewQuery().
		Where("user_id", "=", 1).
		OrderBy("created_at", "asc")
	alerts, err := pgmodels.AlertViewSelect(query)
	require.Nil(t, err)
	assert.Equal(t, 6, len(alerts))
	for _, alertView := range alerts {
		assert.NotEmpty(t, alertView.Content)
		assert.NotEmpty(t, alertView.Subject)
		assert.NotEmpty(t, alertView.Type)
		assert.NotEmpty(t, alertView.UserID)
		assert.NotEmpty(t, alertView.UserEmail)
		assert.NotEmpty(t, alertView.UserName)
	}
}

func TestAlertViewHasBeenRead(t *testing.T) {
	a := &pgmodels.AlertView{}
	assert.False(t, a.HasBeenRead())
	a.ReadAt = time.Now().UTC()
	assert.True(t, a.HasBeenRead())
}
