package admin_api_test

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"os"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	tu "github.com/APTrust/registry/web/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareFileDelete(t *testing.T) {
	tu.InitHTTPTests(t)
	defer db.ForceFixtureReload()

	resp := tu.SysAdminClient.POST("/admin-api/v3/prepare_file_delete/11").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect()
	resp.Status(http.StatusOK)
	workItem := &pgmodels.WorkItem{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), workItem)
	require.Nil(t, err)
	require.NotNil(t, workItem)
	assert.True(t, workItem.ID > 0)

	// Non sys-admin cannot create any events, period.
	tu.Inst1AdminClient.POST("/admin-api/v3/prepare_file_delete/16").
		WithHeader(constants.APIUserHeader, tu.Inst1Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst1UserClient.POST("/admin-api/v3/prepare_file_delete/16").
		WithHeader(constants.APIUserHeader, tu.Inst1User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2AdminClient.POST("/admin-api/v3/prepare_file_delete/16").
		WithHeader(constants.APIUserHeader, tu.Inst2Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2UserClient.POST("/admin-api/v3/prepare_file_delete/16").
		WithHeader(constants.APIUserHeader, tu.Inst2User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
}

func TestPrepareObjectDelete(t *testing.T) {
	tu.InitHTTPTests(t)
	defer db.ForceFixtureReload()

	resp := tu.SysAdminClient.POST("/admin-api/v3/prepare_object_delete/2").
		WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect()
	resp.Status(http.StatusOK)
	workItem := &pgmodels.WorkItem{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), workItem)
	require.Nil(t, err)
	require.NotNil(t, workItem)
	assert.True(t, workItem.ID > 0)

	// Non sys-admin cannot create any events, period.
	tu.Inst1AdminClient.POST("/admin-api/v3/prepare_object_delete/3").
		WithHeader(constants.APIUserHeader, tu.Inst1Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst1UserClient.POST("/admin-api/v3/prepare_object_delete/3").
		WithHeader(constants.APIUserHeader, tu.Inst1User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2AdminClient.POST("/admin-api/v3/prepare_object_delete/3").
		WithHeader(constants.APIUserHeader, tu.Inst2Admin.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
	tu.Inst2UserClient.POST("/admin-api/v3/prepare_object_delete/3").
		WithHeader(constants.APIUserHeader, tu.Inst2User.Email).
		WithHeader(constants.APIKeyHeader, "password").
		Expect().
		Status(http.StatusForbidden)
}

func TestDeletionEndpointSafeguards(t *testing.T) {
	tu.InitHTTPTests(t)
	config := common.Context().Config
	currentEnv := os.Getenv("APT_ENV")
	defer func() {
		config.EnvName = currentEnv
	}()

	unsafeEnvs := []string{
		"production",
		"demo",
		"staging",
	}

	for _, env := range unsafeEnvs {
		config.EnvName = env

		tu.SysAdminClient.POST("/admin-api/v3/prepare_file_delete/11").
			WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
			WithHeader(constants.APIKeyHeader, "password").
			Expect().Status(http.StatusMethodNotAllowed)

		tu.SysAdminClient.POST("/admin-api/v3/prepare_object_delete/2").
			WithHeader(constants.APIUserHeader, tu.SysAdmin.Email).
			WithHeader(constants.APIKeyHeader, "password").
			Expect().Status(http.StatusMethodNotAllowed)
	}
}
