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
		sql := fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, table)
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

// Load fixture data from a single CSV file.
func loadCSVFile(db *pg.DB, table string) error {
	panicOnWrongEnv()
	file := filepath.Join("db", "fixtures", table+".csv")
	data, err := common.LoadRelativeFile(file)
	if err != nil {
		return fmt.Errorf("File read error in %s: %v", file, err)
	}
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("CSV parse error in %s: %v", file, err)
	}

	// BEGIN TRANSACTION
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Close()

	// Insert rows one at a time. We'll commit them all below.
	// This is much faster than committing on each insert.
	placeholders := sqlPlaceholders(len(records[0]))
	var cols []string
	var colString string
	for i, record := range records {
		if i == 0 {
			// First line of csv file has column names.
			cols = record
			colString = strings.Join(cols, ", ")
			continue
		}
		// Insert a single row
		sql := fmt.Sprintf(`insert into "%s" (%s) values (%s)`, table, colString, placeholders)
		fmt.Println(sql, record)
		values := interfaceSlice(record)
		_, err = db.Exec(sql, values...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("SQL insert error, %s line %d: %v", file, i+1, err)
		}
	}

	// COMMIT TRANSACTION
	// and return any error to the caller
	return tx.Commit()
}

// Get the placeholders for a sql query.
func sqlPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func interfaceSlice(s []string) []interface{} {
	iSlice := make([]interface{}, len(s))
	for i, v := range s {
		if v == "NULL" {
			iSlice[i] = nil
		} else {
			iSlice[i] = v
		}
	}
	return iSlice
}

// Our fixtures insert records with ids 1, 2, 3, etc.
// This resets the Postgres ID sequences for the specified table
// so we don't get conflicting IDs when our tests insert new records.
//
// See: https://stackoverflow.com/questions/244243/how-to-reset-postgres-primary-key-sequence-when-it-falls-out-of-sync
func resetSequences(db *pg.DB) error {
	for _, table := range LoadOrder {
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

// Blow up and die if this is run in any environment other than "test"
// or "integration". We call this at every step of the way, in case some
// clever developer ever tries to use or abuse any function in this file.
// Our paranoid level of protection comes from us actually having deleted
// a production database after mistyping a single character in a command.
// Top-notch DevOps saved the day, but paranoid programming would have
// prevented it in the first place.
func panicOnWrongEnv() {
	envName := os.Getenv("APT_ENV")
	if envName != "test" && envName != "integration" {
		panic("Cannot run destructive DB operations outside test and integration environments.")
	}
}
