package db

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
)

var fixturesLoaded = false

// LoadOrder lists the names of tables for which we have fixture data
// (csv files) in the order they should be loaded.
var LoadOrder = []string{
	"roles",
	"institutions",
	"users",
	"roles_users",
	"intellectual_objects",
	"generic_files",
	"checksums",
	"storage_records",
	"premis_events",
	"work_items",
}

// DropOrder lists tables to be dropped, in the order they should be
// dropped so we don't violate foreign key constraints.
var DropOrder = []string{
	"ar_internal_metadata",
	"bulk_delete_jobs",
	"bulk_delete_jobs_emails",
	"bulk_delete_jobs_generic_files",
	"bulk_delete_jobs_institutions",
	"bulk_delete_jobs_intellectual_objects",
	"confirmation_tokens",
	"emails",
	"emails_generic_files",
	"emails_intellectual_objects",
	"emails_premis_events",
	"emails_work_items",
	"old_passwords",
	"schema_migrations",
	"snapshots",
	"usage_samples",
	"work_items",
	"premis_events",
	"storage_records",
	"checksums",
	"generic_files",
	"intellectual_objects",
	"roles_users",
	"users",
	"institutions",
	"roles",
}

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

func dropEverything(db *pg.DB) error {
	panicOnWrongEnv()
	// Drop all tables
	for _, table := range DropOrder {
		sql := fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, table)
		err := runTransaction(db, sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadSchema(db *pg.DB) {
	panicOnWrongEnv()
	// Run the entire schema
}

func loadCSVFiles(db *pg.DB) {
	panicOnWrongEnv()
	// load csv files
}

func loadCSVFile(db *pg.DB, table string) error {
	panicOnWrongEnv()
	file := filepath.Join("db", "fixtures", table+".csv")
	data, err := common.LoadRelativeFile(file)
	if err != nil {
		return err
	}
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// BEGIN TRANSACTION

	//var cols []string
	for i, record := range records {
		if i == 0 {
			// First line of csv file has column names
			// cols = record
		}
		// `insert into (cols...) values (?,?,?...)`, params
	}

	// COMMIT TRANSACTION

	return nil
}

func resetSequence(db *pg.DB, table string) error {
	panicOnWrongEnv()
	//
	// Reset Postgres ID sequences for all tables
	// so we don't get conflicting IDs when our
	// tests insert new records.
	//
	// See: https://stackoverflow.com/questions/244243/how-to-reset-postgres-primary-key-sequence-when-it-falls-out-of-sync
	sql := fmt.Sprintf(`SELECT setval(pg_get_serial_sequence('%s', 'id'), MAX(id)) FROM "%s";`, table, table)
	return runTransaction(db, sql)
}

func runTransaction(db *pg.DB, sql string, params ...interface{}) error {
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		_, err := db.Exec(sql, params)
		return err
	})
}

func panicOnWrongEnv() {
	envName := os.Getenv("APT_ENV")
	if envName != "test" && envName != "integration" {
		panic("Cannot run destructive DB operations outside test and integration environments.")
	}
}
