package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/models"
	"github.com/stretchr/testify/assert"
)

// Bloomsday
var TestDate = time.Date(2021, 6, 16, 10, 24, 16, 0, time.UTC)

const (

	// Test constants for users (from fixture data)
	SysAdmin     = "system@aptrust.org"
	InstAdmin    = "admin@inst1.edu"
	InstUser     = "user@inst1.edu"
	InactiveUser = "inactive@inst1.edu"
	Password     = "password"

	// Institution IDs (from fixture data)
	InstAPTrust = int64(1)
	InstOne     = int64(2)
	InstTwo     = int64(3)
	InstTest    = int64(4)
	InstExample = int64(5)
)

// getUser returns a User object with basic properties set.
// The caller can adjust properties as necessary after the object is created.
// This user has ID zero and is not saved to the DB. By default, user will be
// inst admin at InstOne.
func getUser() (*models.User, error) {
	nonce := time.Now().UnixNano()
	pwd, err := common.EncryptPassword(Password)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Name:              fmt.Sprintf("User %d", nonce),
		Email:             fmt.Sprintf("%d@example.com", nonce),
		EncryptedPassword: pwd,
		InstitutionID:     InstOne,
		Role:              &models.Role{Name: constants.RoleInstAdmin},
	}, nil
}

func TestTypeOf(t *testing.T) {
	user := &models.User{}
	assert.Equal(t, "User", models.TypeOf(user))
}
