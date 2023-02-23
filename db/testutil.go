package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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
	"travis",
}

// LoadOrder lists the names of tables for which we have fixture data
// (csv files) in the order they should be loaded.
var LoadOrder = []string{
	"storage_options",
	"institutions",
	"users",
	"intellectual_objects",
	"generic_files",
	"checksums",
	"storage_records",
	"premis_events",
	"work_items",
	"deletion_requests",
	"deletion_requests_generic_files",
	"deletion_requests_intellectual_objects",
	"alerts",
	"alerts_premis_events",
	"alerts_users",
	"alerts_work_items",
}

// HasNoIDColumn lists tables that have no identity column. Attempting
// to reset the id sequences on these tables will cause an error.
var HasNoIDColumn = []string{
	"roles_users",
	"deletion_requests_generic_files",
	"deletion_requests_intellectual_objects",
	"alerts_premis_events",
	"alerts_users",
	"alerts_work_items",
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
	"alerts_work_items",
	"alerts_users",
	"alerts_premis_events",
	"alerts",
	"deletion_requests_generic_files",
	"deletion_requests_intellectual_objects",
	"deletion_requests",
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
	"storage_options",
	"historical_deposit_stats",
}

var MaterializedViewDropOrder = []string{
	"current_deposit_stats",
	"premis_event_counts",
	"intellectual_object_counts",
	"generic_file_counts",
	"work_item_counts",
}

// LoadFixtures wipes out and resets the test database, loading
// all of the fixtures we use for unit and intergration tests.
// This method effectively runs only once per test run. All calls
// to this function after the first call are no-ops to prevent
// time-consuming reloading of DB for each test.
//
// If you truly want to force a reloading of all fixtures, call
// ForceFixtureReload().
//
// This will panic if run in any APT_ENV other than "test" or "integration".
func LoadFixtures() error {
	panicOnWrongEnv()
	if !fixturesLoaded {
		ctx := common.Context()
		if err := dropEverything(ctx.DB); err != nil {
			ctx.Log.Error().Stack().Err(err).Msg("")
			return err
		}
		if err := loadSchema(ctx.DB); err != nil {
			ctx.Log.Error().Stack().Err(err).Msg("")
			return err
		}
		if err := runMigrations(ctx.DB); err != nil {
			ctx.Log.Error().Stack().Err(err).Msg("")
			return err
		}
		if err := loadCSVFiles(ctx.DB); err != nil {
			ctx.Log.Error().Stack().Err(err).Msg("")
			return err
		}
		if err := resetSequences(ctx.DB); err != nil {
			ctx.Log.Error().Stack().Err(err).Msg("")
			return err
		}
		if err := populateCountsAndStats(ctx.DB); err != nil {
			ctx.Log.Error().Stack().Err(err).Msg("")
			return err
		}
		fixturesLoaded = true
	}
	return nil
}

// ForceFixtureReload forces reloading of all test fixtures.
// This is useful when you want to ensure a clean slate, wiping
// out all records left by prior tests.
func ForceFixtureReload() error {
	fixturesLoaded = false
	return LoadFixtures()
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
	for _, materializedView := range MaterializedViewDropOrder {
		sql := fmt.Sprintf(`drop materialized view if exists "%s"`, materializedView)
		err := runTransaction(db, sql)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reload the entire DB table schema.
// Views will be loaded separately below.
func loadSchema(db *pg.DB) error {
	panicOnWrongEnv()
	file := filepath.Join("db", "schema.sql")
	ddl, err := common.LoadRelativeFile(file)
	if err != nil {
		return fmt.Errorf("File %s: %v", file, err)
	}
	return runTransaction(db, string(ddl))
}

// Run all db migrations
func runMigrations(db *pg.DB) error {
	panicOnWrongEnv()
	dir := filepath.Join(common.ProjectRoot(), "db", "migrations")
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file := info.Name()
		if info.Mode().IsRegular() && strings.HasSuffix(file, ".sql") {
			absPath := filepath.Join(dir, file)
			ddl, err := ioutil.ReadFile(absPath)
			if err != nil {
				return fmt.Errorf("File %s: %v", file, err)
			}
			err = runTransaction(db, string(ddl))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Load all fixture data from CSV files.
func loadCSVFiles(db *pg.DB) error {
	panicOnWrongEnv()
	for _, table := range LoadOrder {
		err := loadCSVFile(db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

// Loads data from a CSV file into a table.
// The CSV files were created with the Postgres COPY command.
func loadCSVFile(db *pg.DB, table string) error {
	panicOnWrongEnv()
	ctx := common.Context()
	file := filepath.Join(common.ProjectRoot(), "db", "fixtures", table+".csv")

	// On Travis, posgres user can't read from Travis' home dir,
	// so we have to copy our csv file to a readable temp dir.
	if ctx.Config.EnvName == "travis" {
		tmpFile := path.Join(os.TempDir(), table+".csv")
		err := common.CopyFile(file, tmpFile, 0666)
		if err != nil {
			return err
		}
		defer os.Remove(tmpFile)
		file = tmpFile
	}

	sql := fmt.Sprintf(`copy "%s" from '%s' csv header`, table, file)
	err := runTransaction(db, sql)
	if err != nil {
		err = fmt.Errorf(`Error executing "%s": %v`, sql, err)
	}
	return err
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

// populate tables and materialized views that cache data from
// expensive queries: counts and deposit stats.
func populateCountsAndStats(db *pg.DB) error {
	// Populate the counts in our *_counts (materialized views).
	_, err := db.Exec("select update_counts();")
	if err != nil {
		return err
	}

	// Populate current deposit stats (materialized view).
	_, err = db.Exec("select update_current_deposit_stats();")
	if err != nil {
		return err
	}

	// Since historical_deposit_stats is a table and not a
	// materialized view, we need to empty it before re-populating it.
	// This table is initially populated by the migration 001_deposit_stats.sql
	// before the fixtures are loaded, so it will be full of zeroes.
	// We want to empty and re-populate it after fixtures are loaded,
	// so it has non-zero values.
	_, err = db.Exec("delete from historical_deposit_stats")
	if err != nil {
		return err
	}

	// Populate all historical deposit stats (data goes into a table).
	_, err = db.Exec("select populate_all_historical_deposit_stats();")
	return err
}

// Blow up and die if this is run in any environment other than "test",
// "integration", "travis", or "dev". We call this at every step of
// the way, in case some clever developer ever tries to use or abuse
// any function in this file.
//
// Our paranoid level of protection comes from us actually having deleted
// a production database after mistyping a single character in a rake
// command. Top-notch DevOps saved the day, but paranoid programming
// would have prevented it in the first place.
func panicOnWrongEnv() {
	ctx := common.Context()
	envName := ctx.Config.EnvName
	if !slice.Contains(SafeEnvironments, envName) {
		panic("Cannot run destructive DB operations outside dev, integration, test, travis environments.")
	}
	// Be extra safe. Why? Because some jackass rake task once
	// deleted our entire production DB.
	//
	// Don't like the kludgy, ugly code below? See how much you
	// like restoring 250GB of deleted data.
	if ctx.Config.DB.Host != "localhost" {
		panic("Cannot run destructive DB operations against external servers.")
	}
	if !strings.HasSuffix(ctx.Config.DB.Name, "_development") &&
		!strings.HasSuffix(ctx.Config.DB.Name, "_integration") &&
		!strings.HasSuffix(ctx.Config.DB.Name, "_test") &&
		!strings.HasSuffix(ctx.Config.DB.Name, "_travis") {
		panic("Cannot run destructive DB operations against non dev, integration, test, or travis DB.")
	}
}
