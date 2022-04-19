package admin_api_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: When using the httpexpect client, every time you call Expect()
// on a POST or other request, it re-sends the request to the server
// and it screws up the Request URL. First request, is /admin-api/v3/events/create,
// second is /admin-api/v3/events/create/admin-api/v3/events/create, third is
// /admin-api/v3/events/create/admin-api/v3/events/create/admin-api/v3/events/create,
// etc.
//
// These subsequent requests cause errors in the gin engine because they don't
// match any known routes.
func TestPremisEventCreate(t *testing.T) {
	tu.InitHTTPTests(t)

	gf, err := pgmodels.GenericFileByID(21)
	require.Nil(t, err)
	require.NotNil(t, gf)
	event := pgmodels.RandomPremisEvent(constants.EventIngestion)
	event.GenericFileID = gf.ID
	event.IntellectualObjectID = gf.IntellectualObjectID

	jsonData, err := json.Marshal(event)
	require.Nil(t, err)
	require.NotEmpty(t, jsonData)

	resp := tu.SysAdminClient.POST("/admin-api/v3/events/create").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithBytes(jsonData).
		Expect()
	savedEvent := &pgmodels.PremisEvent{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), savedEvent)
	require.Nil(t, err)
	require.NotNil(t, savedEvent)
	assert.True(t, savedEvent.ID > 0)

	// Non sys-admin cannot create any events, period.
	tu.Inst1AdminClient.POST("/admin-api/v3/events/create").
		WithHeader(constants.APIUserHeader, tu.Inst1Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst1UserClient.POST("/admin-api/v3/events/create").
		WithHeader(constants.APIUserHeader, tu.Inst1User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2AdminClient.POST("/admin-api/v3/events/create").
		WithHeader(constants.APIUserHeader, tu.Inst2Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2UserClient.POST("/admin-api/v3/events/create").
		WithHeader(constants.APIUserHeader, tu.Inst2User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)
}

func TestPremisEventFixityCheck(t *testing.T) {
	tu.InitHTTPTests(t)

	gf, err := pgmodels.GenericFileByID(21)
	require.Nil(t, err)
	require.NotNil(t, gf)

	now := time.Now().UTC()
	assert.NotEqual(t, now, gf.LastFixityCheck)

	event := pgmodels.RandomPremisEvent(constants.EventFixityCheck)
	event.GenericFileID = gf.ID
	event.IntellectualObjectID = gf.IntellectualObjectID
	event.DateTime = now

	jsonData, err := json.Marshal(event)
	require.Nil(t, err)
	require.NotEmpty(t, jsonData)

	resp := tu.SysAdminClient.POST("/admin-api/v3/events/create").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		WithBytes(jsonData).
		Expect()
	savedEvent := &pgmodels.PremisEvent{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), savedEvent)
	require.Nil(t, err)
	require.NotNil(t, savedEvent)
	assert.True(t, savedEvent.ID > 0)

	// Make sure LastFixityCheck was upated
	gf, err = pgmodels.GenericFileByID(21)
	require.Nil(t, err)
	require.NotNil(t, gf)
	assert.Equal(t, now, gf.LastFixityCheck)
}
