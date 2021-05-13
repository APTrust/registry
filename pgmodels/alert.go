package pgmodels

import (
	"time"

	"github.com/go-pg/pg/v10/orm"
)

type Alert struct {
	ID                int64            `json:"id"`
	InstitutionID     int64            `json:"institution_id"`
	Type              string           `json:"type"`
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
	SentAt        time.Time
	ReadAt        time.Time
}

type AlertsUsers struct {
	AlertID int64
	UserID  int64
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
