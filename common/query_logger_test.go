package common_test

import (
	"context"
	"strings"
	"testing"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryLogger(t *testing.T) {
	var sb = strings.Builder{}
	logger := common.NewQueryLogger(zerolog.New(&sb))
	require.NotNil(t, logger)

	event := &pg.QueryEvent{}
	_, err := logger.BeforeQuery(context.Background(), event)
	require.Nil(t, err)
	assert.Equal(t, `{"level":"debug","message":"Starting SQL: "}`+"\n", sb.String())

	err = logger.AfterQuery(context.Background(), event)
	require.Nil(t, err)

}
