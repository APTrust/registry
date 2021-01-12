package db

import (
	"os"
	//"github.com/APTrust/registry/common"
)

var fixturesLoaded = false

func LoadFixtures() {
	panicOnWrongEnv()
	if !fixturesLoaded {
		//ctx := common.Context()

		// dropEverything()
		// loadSchema()
		// loadCSVFiles()
		// resetSequences()

		fixturesLoaded = true
	}
}

func dropEverything() {
	// Drop all tables
}

func loadSchema() {
	// Run the entire schema
}

func loadCSVFiles() {
	// load csv files
}

func resetSequences() {
	//
	// Reset Postgres ID sequences for all tables
	// so we don't get conflicting IDs when our
	// tests insert new records.
	//
	// SELECT setval(pg_get_serial_sequence('t1', 'id'), coalesce(max(id),0) + 1, false) FROM t1;
}

func panicOnWrongEnv() {
	envName := os.Getenv("APT_ENV")
	if envName != "test" && envName != "integration" {
		panic("Cannot run destructive DB operations outside test and integration environments.")
	}
}
