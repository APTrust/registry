package constants_test

import (
	"testing"

	"github.com/APTrust/registry/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTopicFor(t *testing.T) {
	topic, err := constants.TopicFor(constants.ActionDelete, "")
	require.Nil(t, err)
	assert.Equal(t, constants.TopicDelete, topic)

	topic, err = constants.TopicFor(constants.ActionRestoreFile, "")
	require.Nil(t, err)
	assert.Equal(t, constants.TopicFileRestore, topic)

	topic, err = constants.TopicFor(constants.ActionGlacierRestore, "")
	require.Nil(t, err)
	assert.Equal(t, constants.TopicGlacierRestore, topic)

	topic, err = constants.TopicFor(constants.ActionRestoreObject, "")
	require.Nil(t, err)
	assert.Equal(t, constants.TopicObjectRestore, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageReceive)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestPreFetch, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageValidate)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestValidation, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageReingestCheck)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestReingestCheck, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageCopyToStaging)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestStaging, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageFormatIdentification)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestFormatIdentification, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageStore)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestStorage, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageStorageValidation)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestStorageValidation, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageRecord)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestRecord, topic)

	topic, err = constants.TopicFor(constants.ActionIngest, constants.StageCleanup)
	require.Nil(t, err)
	assert.Equal(t, constants.IngestCleanup, topic)

	topic, err = constants.TopicFor("invalid action", "invalid stage")
	require.NotNil(t, err)
	assert.Empty(t, topic)
}
