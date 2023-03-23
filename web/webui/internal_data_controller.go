package webui

import (
	"net/http"

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

	migrations, err := pgmodels.SchemaMigrationSelect(pgmodels.NewQuery().OrderBy("version", "desc"))
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["migrations"] = migrations

	c.HTML(http.StatusOK, "internal_metadata/index.html", req.TemplateData)
}
