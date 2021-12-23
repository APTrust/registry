package admin_api_test

import (
	"encoding/json"
	"fmt"
	//"net/http"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	tu "github.com/APTrust/registry/web/testutil"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// START HERE
//
// TODO: The following test is failing because middleware.ResourceAuthorization
//       can't determine the resource type. It's getting the CSRF handler as
//       the function name, which shouldn't happen, esp. because CSRF should
//       skip API requests. The handler name should be PremisEventCreate,
//       and the associated resource type is PremisEvent. So WTF?
//
//       Happy New Year!

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

	resp := tu.SysAdminClient.POST("/admin-api/v3/events/create").WithBytes(jsonData)

	//fmt.Println(resp.Expect())
	fmt.Println(string(resp.Expect().Body().Raw()))

	//resp.Expect().Status(http.StatusCreated)

	// savedEvent := &pgmodels.PremisEvent{}
	// err = json.Unmarshal([]byte(resp.Expect().Body().Raw()), savedEvent)
	// require.Nil(t, err)
	// require.NotNil(t, savedEvent)
	// assert.True(t, savedEvent.ID > 0)

	// ------------------------

	// // Non sys-admin cannot create any events, period.
	// instIds := []int64{1, 2, 3, 4, 5}
	// for _, id := range instIds {
	// 	tu.Inst1AdminClient.GET("/admin-api/v3/institutions/show/{id}", id).
	// 		Expect().
	// 		Status(http.StatusForbidden)
	// 	tu.Inst1UserClient.GET("/admin-api/v3/institutions/show/{id}", id).
	// 		Expect().
	// 		Status(http.StatusForbidden)
	// 	tu.Inst2AdminClient.GET("/admin-api/v3/institutions/show/{id}", id).
	// 		Expect().
	// 		Status(http.StatusForbidden)
	// 	tu.Inst2UserClient.GET("/admin-api/v3/institutions/show/{id}", id).
	// 		Expect().
	// 		Status(http.StatusForbidden)
	// }
}
