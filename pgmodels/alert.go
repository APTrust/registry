package pgmodels

import (
	"context"
	"strings"
	"time"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/constants"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/stew/slice"
)

const (
	ErrAlertInstitutionID = "InstitutionID is required."
	ErrAlertType          = "Alert type is missing or invalid."
	ErrAlertContent       = "Alert content cannot be empty."
)

type Alert struct {
	ID                int64            `json:"id"`
	InstitutionID     int64            `json:"institution_id"`
	Type              string           `json:"type"`
	Subject           string           `json:"subject"`
	Content           string           `json:"content"`
	DeletionRequestID int64            `json:"deletion_request_id"`
	CreatedAt         time.Time        `json:"created_at"`
	DeletionRequest   *DeletionRequest `json:"-" pg:"rel:has-one"`
	PremisEvents      []*PremisEvent   `json:"premis_events" pg:"many2many:alerts_premis_events"`
	Users             []*User          `json:"users" pg:"many2many:alerts_users"`
	WorkItems         []*WorkItem      `json:"work_items" pg:"many2many:alerts_work_items"`
}

type AlertsPremisEvents struct {
	AlertID       int64
	PremisEventID int64
}

type AlertsUsers struct {
	AlertID int64
	UserID  int64
	SentAt  time.Time
	ReadAt  time.Time
}

type AlertsWorkItems struct {
	AlertID    int64
	WorkItemID int64
}

// init does some setup work so go-pg can recognize many-to-many
// relations. Go automatically calls this function once when package
// is imported.
func init() {
	orm.RegisterTable((*AlertsPremisEvents)(nil))
	orm.RegisterTable((*AlertsUsers)(nil))
	orm.RegisterTable((*AlertsWorkItems)(nil))
}

// AlertByID returns the alert with the specified id.
// Returns pg.ErrNoRows if there is no match.
func AlertByID(id int64) (*Alert, error) {
	query := NewQuery().Where(`"alert"."id"`, "=", id).Relations("DeletionRequest", "PremisEvents", "Users", "WorkItems")
	return AlertGet(query)
}

// AlertGet returns the first alert matching the query.
func AlertGet(query *Query) (*Alert, error) {
	var alert Alert
	err := query.Select(&alert)
	return &alert, err
}

// AlertSelect returns all alerts matching the query.
func AlertSelect(query *Query) ([]*Alert, error) {
	var alerts []*Alert
	err := query.Select(&alerts)
	return alerts, err
}

func (alert *Alert) GetID() int64 {
	return alert.ID
}

// Save saves this alert to the database. This will peform an insert
// if Alert.ID is zero. Otherwise, it updates. It also saves all of
// the many-to-many relations (PremisEvents, Users, and WorkItems), though
// note that on update it does not delete any of these relations. We don't
// have a use case for that yet, since alerts are generally created and never
// updated.
func (alert *Alert) Save() error {
	registryContext := common.Context()
	db := registryContext.DB
	return db.RunInTransaction(db.Context(), func(*pg.Tx) error {
		var err error
		if alert.ID == 0 {
			_, err = db.Model(alert).Insert()
		} else {
			_, err = db.Model(alert).WherePK().Update()
		}
		if err != nil {
			registryContext.Log.Error().Msgf("Transaction failed. Model: %v. Error: %v", alert, err)
		}
		return alert.saveRelations(db)
	})
}

// This is run inside the Save transaction.
func (alert *Alert) saveRelations(db *pg.DB) error {
	err := alert.saveEvents(db)
	if err != nil {
		return err
	}
	err = alert.saveWorkItems(db)
	if err != nil {
		return err
	}
	err = alert.saveUsers(db)
	return err
}

func (alert *Alert) saveEvents(db *pg.DB) error {
	sql := "insert into alerts_premis_events (alert_id, premis_event_id) values (?, ?) on conflict do nothing"
	for _, event := range alert.PremisEvents {
		_, err := db.Exec(sql, alert.ID, event.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (alert *Alert) saveUsers(db *pg.DB) error {
	sql := "insert into alerts_users (alert_id, user_id, sent_at, read_at) values (?, ?, ?, ?) on conflict do nothing"
	for _, user := range alert.Users {
		_, err := db.Exec(sql, alert.ID, user.ID, nil, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (alert *Alert) saveWorkItems(db *pg.DB) error {
	sql := "insert into alerts_work_items (alert_id, work_item_id) values (?, ?) on conflict do nothing"
	for _, item := range alert.WorkItems {
		_, err := db.Exec(sql, alert.ID, item.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// The following statements have no effect other than to force a compile-time
// check that ensures our Alert model properly implements these hook
// interfaces.
var (
	_ pg.BeforeInsertHook = (*Alert)(nil)
	_ pg.BeforeUpdateHook = (*Alert)(nil)
)

// BeforeInsert sets timestamps and bucket names on creation.
func (alert *Alert) BeforeInsert(c context.Context) (context.Context, error) {
	now := time.Now().UTC()
	alert.CreatedAt = now

	err := alert.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// BeforeUpdate sets the UpdatedAt timestamp.
func (alert *Alert) BeforeUpdate(c context.Context) (context.Context, error) {
	err := alert.Validate()
	if err == nil {
		return c, nil
	}
	return c, err
}

// Validate validates the model. This is called automatically on insert
// and update.
func (alert *Alert) Validate() *common.ValidationError {
	errors := make(map[string]string)

	if alert.InstitutionID < 1 {
		errors["InstitutionID"] = ErrAlertInstitutionID
	}
	if !slice.Contains(constants.AlertTypes, alert.Type) {
		errors["Type"] = ErrAlertType
	}
	if strings.TrimSpace(alert.Content) == "" {
		errors["Content"] = ErrAlertContent
	}

	if len(errors) > 0 {
		return &common.ValidationError{Errors: errors}
	}
	return nil
}
