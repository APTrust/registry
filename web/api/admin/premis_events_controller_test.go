package admin_api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"os"
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
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
//       The root of the problem is that the test client is requesting
//               /admin-api/v3/events/create/
//       but the server is seeing
//               /admin-api/v3/events/create/admin-api/v3/events/create
//
//       Gin is doing some redirects on its own, but I can't figure out
//       where. In fact, it processes requests to /admin-api/v3/events/create/
//       three times, with the final request haveing the the route as
//       /admin-api/v3/events/create/admin-api/v3/events/create/admin-api/v3/events/create/
//
//       Why??
//
//       Dunno, but tests for non-admin users report a 307 redirect.
//       Just have to figure out where that's coming from and why.
//
//       Also note that the APTrust session cookie is being set twice.
//       All redirects force all the middleware to fire again.

func TestPremisEventCreate(t *testing.T) {
	//os.Setenv("APT_ENV", "test")
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

	fmt.Println("Client request path:", resp.Expect().Raw().Request.URL.Path)

	//fmt.Println(resp.Expect())
	fmt.Println(string(resp.Expect().Body().Raw()))
	fmt.Println(resp.Expect().Raw().StatusCode)

	//resp.Expect().Status(http.StatusCreated)

	savedEvent := &pgmodels.PremisEvent{}
	err = json.Unmarshal([]byte(resp.Expect().Body().Raw()), savedEvent)
	require.Nil(t, err)
	require.NotNil(t, savedEvent)
	assert.True(t, savedEvent.ID > 0)

	// Non sys-admin cannot create any events, period.
	tu.Inst1AdminClient.POST("/admin-api/v3/events/create").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst1UserClient.POST("/admin-api/v3/events/create").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2AdminClient.POST("/admin-api/v3/events/create").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2UserClient.POST("/admin-api/v3/events/create").
		WithBytes(jsonData).
		Expect().
		Status(http.StatusForbidden)

}
