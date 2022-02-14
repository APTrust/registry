package common_api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/APTrust/registry/pgmodels"
	"github.com/APTrust/registry/web/api"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPremisEventShow(t *testing.T) {
	tu.InitHTTPTests(t)

	// In fixture data, event #2 belongs to inst1
	event, err := pgmodels.PremisEventByID(2)
	require.Nil(t, err)
	require.NotNil(t, event)

	// Sysadmin can read any event
	resp := tu.SysAdminClient.GET("/member-api/v3/events/show/{id}", event.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.PremisEvent{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, event.ID, record.ID)
	assert.Equal(t, event.InstitutionID, record.InstitutionID)

	// Make sure we can get this event by identifier as well.
	resp = tu.SysAdminClient.GET("/member-api/v3/events/show/{id}", event.Identifier).Expect().Status(http.StatusOK)
	record = &pgmodels.PremisEvent{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, event.ID, record.ID)
	assert.Equal(t, event.Identifier, record.Identifier)

	// Inst admin can read event from own inst
	resp = tu.Inst1AdminClient.GET("/member-api/v3/events/show/{id}", event.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, event.ID, record.ID)
	assert.Equal(t, event.InstitutionID, record.InstitutionID)

	// Inst admin CANNOT read event from other institution
	tu.Inst2AdminClient.GET("/member-api/v3/events/show/{id}", event.ID).
		Expect().
		Status(http.StatusForbidden)

	// Inst user can read event from own inst
	resp = tu.Inst1UserClient.GET("/member-api/v3/events/show/{id}", event.ID).Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, event.ID, record.ID)
	assert.Equal(t, event.InstitutionID, record.InstitutionID)

	// Inst user CANNOT read event from other institution
	tu.Inst2UserClient.GET("/member-api/v3/events/show/{id}", event.ID).
		Expect().
		Status(http.StatusForbidden)

}

func TestPremisEventIndex(t *testing.T) {
	tu.InitHTTPTests(t)

	// Sys Admin should see all events and filters
	resp := tu.SysAdminClient.GET("/member-api/v3/events").
		WithQuery("page", 2).
		WithQuery("per_page", 5).
		Expect().Status(http.StatusOK)

	list := api.PremisEventViewList{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 54, list.Count)
	assert.Equal(t, "/member-api/v3/events?page=3&per_page=5", list.Next)
	assert.Equal(t, "/member-api/v3/events?page=1&per_page=5", list.Previous)
	assert.Equal(t, tu.Inst2User.InstitutionID, list.Results[0].InstitutionID)

	// Test some filters. This object has 1 deleted, 4 active events.
	resp = tu.SysAdminClient.GET("/member-api/v3/events").
		WithQuery("intellectual_object_id", 3).
		Expect().Status(http.StatusOK)

	list = api.PremisEventViewList{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 7, list.Count)
	assert.Equal(t, 7, len(list.Results))
	for _, event := range list.Results {
		assert.Equal(t, int64(3), event.IntellectualObjectID)
		assert.Equal(t, "institution1.edu/glass", event.IntellectualObjectIdentifier)
		assert.Equal(t, int64(2), event.InstitutionID)
	}

	// Inst admin should see only his own institution's events.
	resp = tu.Inst1AdminClient.GET("/member-api/v3/events").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 27, list.Count)
	assert.Equal(t, "/member-api/v3/events?page=2&per_page=20", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 20, len(list.Results))
	for _, event := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, event.InstitutionID)
	}

	// Inst admin cannot see events belonging to other insitutions.
	tu.Inst2AdminClient.GET("/member-api/v3/events").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

	// Inst user should see only his own institution's events.
	resp = tu.Inst1UserClient.GET("/member-api/v3/events").
		Expect().Status(http.StatusOK)
	err = json.Unmarshal([]byte(resp.Body().Raw()), &list)
	require.Nil(t, err)
	assert.Equal(t, 27, list.Count)
	assert.Equal(t, "/member-api/v3/events?page=2&per_page=20", list.Next)
	assert.Equal(t, "", list.Previous)
	assert.Equal(t, 20, len(list.Results))
	for _, event := range list.Results {
		assert.Equal(t, tu.Inst1User.InstitutionID, event.InstitutionID)
	}

	// Inst user cannot see other institution's events.
	tu.Inst2UserClient.GET("/member-api/v3/events").
		WithQuery("institution_id", tu.Inst1Admin.InstitutionID).
		Expect().Status(http.StatusForbidden)

}
