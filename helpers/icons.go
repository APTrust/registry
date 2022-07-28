package helpers

import (
	"github.com/APTrust/registry/constants"
)

// IconMissing is the icon we show for items that are not in the IconMap.
// This appears as a question mark. If you see it on any web
// page, you should know to add an appropriate icon to the IconMap.
var IconMissing = `help_outline`

// IconMap maps strings to Material icons.
var IconMap = map[string]string{

	// Premis Event Icons.
	// Note that we use only 7 or so event types,
	// so we don't define an icon for every type.
	constants.EventAccessAssignment:     `admin_panel_settings`,
	constants.EventCreate:               `add_circle_outline`,
	constants.EventDeletion:             `delete_forever`,
	constants.EventDigestCalculation:    `description`,
	constants.EventFixityCheck:          `fingerprint`,
	constants.EventIdentifierAssignment: `search`,
	constants.EventIngestion:            `file_upload`,
	constants.EventReplication:          `library_books`,
}

// BadgeClassMap maps work item status and other constant values
// to css badge classes.
var BadgeClassMap = map[string]string{
	constants.StatusCancelled: "is-cancelled",
	constants.StatusFailed:    "is-failed",
	constants.StatusPending:   "is-pending",
	constants.StatusStarted:   "is-started",
	constants.StatusSuccess:   "is-success",
	constants.StatusSuspended: "is-suspended",
}
