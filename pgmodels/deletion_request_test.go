package pgmodels_test

import (
	"testing"

	//"github.com/APTrust/registry/constants"
	//"github.com/APTrust/registry/db"
	"github.com/APTrust/registry/pgmodels"
	//"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)

const (
	ConfTokenPlain       = "ConfirmationToken"
	CancelTokenPlain     = "CancelToken"
	ConfTokenEncrypted   = "$2a$10$TK8s1XnmWulSUdze8GN5uOgGmDDsnndQKF5/Rz1j0xaHT7AwXRVma"
	CancelTokenEncrypted = "$2a$10$xwxTFn.k1TbfbNSW3/udduwtjwo7nQSBlIlARHvTXADAhCfQtZt46"
)

var request = &pgmodels.DeletionRequest{}

func TestDeletionRequestByID(t *testing.T) {

}
