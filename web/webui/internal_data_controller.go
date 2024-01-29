package webui

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// InternalDataIndex returns data from the schema_migrations and
// ar_internal_metadata tables.
// GET /internal_metadata
func InternalMetadataIndex(c *gin.Context) {
	req := NewRequest(c)
	internalMetaData, err := pgmodels.InternalMetadataSelect(pgmodels.NewQuery().OrderBy("key", "asc"))
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["internalMetadata"] = internalMetaData

	// Filter out migrations with null started_at timestamps.
	// Those are legacy Rails migrations, and there are a lot of them.
	migrations, err := pgmodels.SchemaMigrationSelect(pgmodels.NewQuery().IsNotNull("started_at").OrderBy("version", "desc"))
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["migrations"] = migrations
	req.TemplateData["envName"] = common.Context().Config.EnvName
	req.TemplateData["redisUrl"] = common.Context().Config.Redis.URL
	req.TemplateData["dbName"] = fmt.Sprintf("%s@%s", common.Context().Config.DB.Name, common.Context().Config.DB.Host)

	c.HTML(http.StatusOK, "internal_metadata/index.html", req.TemplateData)
}
