package main

import (
	"fmt"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/db"
)

// Load db fixtures for integration tests.
// We extracted this into a separately compiled app because we don't
// want this destructive operation to be accessible on any of our
// live systems.
//
// Note that db.LoadFixtures has internal safeguards against running
// in any environment other than a handful of known dev/test environments.
func main() {
	ctx := common.Context()
	err := db.LoadFixtures()
	if err != nil {
		fmt.Println("Error loading fixtures:", err.Error())
	} else {
		fmt.Println("Loaded fixtures for environment", ctx.Config.EnvName)
	}
}
