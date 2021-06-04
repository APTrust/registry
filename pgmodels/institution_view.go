package pgmodels

import (
	"time"
)

// InstitutionView contains information about an institution and its
// parent (if it has one). This view simplifies both search and display
// for institution management.
type InstitutionView struct {
	tableName           struct{}  `pg:"institutions_view"`
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Identifier          string    `json:"identifier"`
	State               string    `json:"state"`
	Type                string    `json:"type"`
	DeactivatedAt       time.Time `json:"deactivated_at"`
	OTPEnabled          bool      `json:"otp_enabled"`
	EnableSpotRestore   bool      `json:"enable_spot_restore"`
	ReceivingBucket     string    `json:"receiving_bucket"`
	RestoreBucket       string    `json:"restore_bucket"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	ParentId            int64     `json:"parent_id"`
	ParentName          string    `json:"parent_name"`
	ParentIdentifier    string    `json:"parent_identifier"`
	ParentState         string    `json:"parent_state"`
	ParentDeactivatedAt time.Time `json:"parent_deactivated_at"`
}

// InstitutionViewByID returns the InstitutionView record
// with the specified id.  Returns pg.ErrNoRows if there is no match.
func InstitutionViewByID(id int64) (*InstitutionView, error) {
	query := NewQuery().Where("id", "=", id)
	return InstitutionViewGet(query)
}

// InstitutionViewByEmail returns the InstitutionView record with the
// specified email address. Returns pg.ErrNoRows if there is no match.
func InstitutionViewByEmail(email string) (*InstitutionView, error) {
	query := NewQuery().Where("email", "=", email)
	return InstitutionViewGet(query)
}

// InstitutionViewSelect returns all InstitutionView records matching
// the query.
func InstitutionViewSelect(query *Query) ([]*InstitutionView, error) {
	var institutions []*InstitutionView
	err := query.Select(&institutions)
	return institutions, err
}

// InstitutionViewGet returns the first user view record matching the query.
func InstitutionViewGet(query *Query) (*InstitutionView, error) {
	var institution InstitutionView
	err := query.Select(&institution)
	return &institution, err
}
