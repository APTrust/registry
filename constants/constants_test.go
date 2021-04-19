package constants_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/stretchr/testify/assert"
)

func TestTopicFor(t *testing.T) {
	assert.Equal(t, constants.TopicDelete, constants.TopicFor(constants.ActionDelete, ""))
	assert.Equal(t, constants.TopicFileRestore, constants.TopicFor(constants.ActionRestoreFile, ""))
	assert.Equal(t, constants.TopicGlacierRestore, constants.TopicFor(constants.ActionGlacierRestore, ""))
	assert.Equal(t, constants.TopicObjectRestore, constants.TopicFor(constants.ActionRestoreObject, ""))

	assert.Equal(t, constants.IngestPreFetch, constants.TopicFor(constants.ActionIngest, constants.StageReceive))
	assert.Equal(t, constants.IngestValidation, constants.TopicFor(constants.ActionIngest, constants.StageValidate))
	assert.Equal(t, constants.IngestReingestCheck, constants.TopicFor(constants.ActionIngest, constants.StageReingestCheck))
	assert.Equal(t, constants.IngestStaging, constants.TopicFor(constants.ActionIngest, constants.StageCopyToStaging))
	assert.Equal(t, constants.IngestFormatIdentification, constants.TopicFor(constants.ActionIngest, constants.StageFormatIdentification))
	assert.Equal(t, constants.IngestStorage, constants.TopicFor(constants.ActionIngest, constants.StageStore))
	assert.Equal(t, constants.IngestStorageValidation, constants.TopicFor(constants.ActionIngest, constants.StageStorageValidation))
	assert.Equal(t, constants.IngestRecord, constants.TopicFor(constants.ActionIngest, constants.StageRecord))
	assert.Equal(t, constants.IngestCleanup, constants.TopicFor(constants.ActionIngest, constants.StageCleanup))
}
