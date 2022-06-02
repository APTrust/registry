package admin_api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/APTrust/registry/pgmodels"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletionRequestShow(t *testing.T) {
	tu.InitHTTPTests(t)

	deletion, err := pgmodels.DeletionRequestByID(2)
	require.Nil(t, err)
	require.NotNil(t, deletion)

	// Sysadmin can read any deletion
	resp := tu.SysAdminClient.GET("/admin-api/v3/deletions/show/{id}", deletion.ID).Expect().Status(http.StatusOK)
	record := &pgmodels.DeletionRequest{}
	err = json.Unmarshal([]byte(resp.Body().Raw()), record)
	require.Nil(t, err)
	assert.Equal(t, deletion.ID, record.ID)
	assert.Equal(t, deletion.InstitutionID, record.InstitutionID)

	// Others cannot read this endpoint. They have to use the member api
	tu.Inst1AdminClient.GET("/admin-api/v3/deletions/show/{id}", deletion.ID).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2AdminClient.GET("/admin-api/v3/deletions/show/{id}", deletion.ID).
		Expect().
		Status(http.StatusForbidden)
	tu.Inst1UserClient.GET("/admin-api/v3/deletions/show/{id}", deletion.ID).
		Expect().
		Status(http.StatusForbidden)
}
