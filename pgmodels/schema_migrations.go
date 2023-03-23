package pgmodels

import "time"

// SchemaMigration represents a schema migration record from
// the database. These tell us which migrations have been run.
// This model is read-only.
type SchemaMigration struct {
	Version    string
	StartedAt  time.Time
	FinishedAt time.Time
}

// SchemaMigrationSelect returns all records matching the query.
func SchemaMigrationSelect(query *Query) ([]*SchemaMigration, error) {
	var records []*SchemaMigration
	err := query.Select(&records)
	return records, err
}
