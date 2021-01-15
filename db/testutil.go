package db

import (
	// "encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/APTrust/registry/common"
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/stew/slice"
)

var fixturesLoaded = false

// SafeEnvironments lists which APT_ENV environments are safe for data loading.
// Since data loading DELETES THE ENTIRE DB before reloading fixtures, we want
// this to run only on local dev machines and test/CI systems.
var SafeEnvironments = []string{
	"dev",
	"integration",
	"test",
}

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

// HasNoIDColumn lists tables that have no identity column. Attempting
// to reset the id sequences on these tables will cause an error.
var HasNoIDColumn = []string{
	"roles_users",
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

// LoadFixtures wipes out and resets the test database, loading
// all of the fixtures we use for unit and intergration tests.
// This method effectively runs only once per test run. All calls
// to this function after the first call are no-ops.
//
// This will panic if run in any APT_ENV other than "test" or "integration".
func LoadFixtures() error {
	panicOnWrongEnv()
	if !fixturesLoaded {
		ctx := common.Context()
		if err := dropEverything(ctx.DB); err != nil {
			return err
		}
		if err := loadSchema(ctx.DB); err != nil {
			return err
		}
		if err := loadCSVFiles(ctx.DB); err != nil {
			return err
		}
		if err := resetSequences(ctx.DB); err != nil {
			return err
		}
		fixturesLoaded = true
	}
	return nil
}

// Drop all tables in the DB.
func dropEverything(db *pg.DB) error {
	panicOnWrongEnv()
	// Drop all tables
	for _, table := range DropOrder {
		sql := fmt.Sprintf(`drop table if exists "%s" cascade`, table)
		err := runTransaction(db, sql)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reload the entire DB schema
func loadSchema(db *pg.DB) error {
	panicOnWrongEnv()
	file := filepath.Join("db", "schema.sql")
	ddl, err := common.LoadRelativeFile(file)
	if err != nil {
		return fmt.Errorf("File %s: %v", file, err)
	}
	return runTransaction(db, string(ddl))
}

// Load all fixture data from CSV files.
func loadCSVFiles(db *pg.DB) error {
	panicOnWrongEnv()
	for _, table := range LoadOrder {
		if err := loadCSVFile(db, table); err != nil {
			return err
		}
	}
	return nil
}

// Loads data from a CSV file into a table.
// The CSV files were created with the Postgres COPY command.
func loadCSVFile(db *pg.DB, table string) error {
	panicOnWrongEnv()
	file := filepath.Join(common.ProjectRoot(), "db", "fixtures", table+".csv")
	sql := fmt.Sprintf(`copy "%s" from '%s' csv header`, table, file)
	err := runTransaction(db, sql)
	if err != nil {
		fmt.Println(sql)
		err = fmt.Errorf(`Error executing "%s": %v`, sql, err)
	}
	return err
}

// Get the placeholders for a sql query.
func sqlPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

// Our fixtures include explicit IDs because other fixtures must refer
// to those IDs in foreign key columns. For example, the intellectual_objects
// fixtures refer to institution ids 1, 2, 3, etc. They MUST do this for
// tests to be reliable.
//
// When we import explicit IDs, we are doing inserts without hitting the
// serial auto-incrementer. The result is that when our tests insert new
// records in tables we've imported, Postgres will start generating IDs
// at 1. But ID 1 already exists, from the imported data.
//
// Here's a sample error:
//
// ERROR: duplicate key value violates unique constraint "<table>_pkey"
// Detail: Key (id)=(1) already exists.
//
// To avoid this, the resetSequence code below tells Postgres to start the
// ID sequence at the highest existing ID in the table, instead of at 1.
//
func resetSequences(db *pg.DB) error {
	for _, table := range LoadOrder {
		if slice.ContainsString(HasNoIDColumn, table) {
			continue
		}
		if err := resetSequence(db, table); err != nil {
			return err
		}
	}
	return nil
}

// Reset the ID sequence on a single table.
func resetSequence(db *pg.DB, table string) error {
	panicOnWrongEnv()
	sql := fmt.Sprintf(`SELECT setval(pg_get_serial_sequence('%s', 'id'), MAX(id)) FROM "%s";`, table, table)
	return runTransaction(db, sql)
}

// This runs a sql command in a transaction, returning an error if one occurs.
func runTransaction(db *pg.DB, sql string, params ...interface{}) error {
	return db.RunInTransaction(db.Context(), func(tx *pg.Tx) error {
		_, err := db.Exec(sql, params)
		return err
	})
}

// Blow up and die if this is run in any environment other than "test",
// "integration", or "dev". We call this at every step of the way, in case
// some clever developer ever tries to use or abuse any function in this file.
// Our paranoid level of protection comes from us actually having deleted
// a production database after mistyping a single character in a command.
// Top-notch DevOps saved the day, but paranoid programming would have
// prevented it in the first place.
func panicOnWrongEnv() {
	envName := os.Getenv("APT_ENV")
	if !slice.Contains(SafeEnvironments, envName) {
		panic("Cannot run destructive DB operations outside dev, integration, test environments.")
	}
}
