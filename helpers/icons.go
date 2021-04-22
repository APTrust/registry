package helpers

import (
	"github.com/APTrust/registry/constants"
)

// IconMissing is the icon we show for items that are not in the IconMap.
// This appears as an upside-down question mark. If you see it on any web
// page, you should know to add an appropriate icon to the IconMap.
var IconMissing = `<i class="far fa-question-circle fa-rotate-180"></i>`

// IconMap maps strings to FontAwesome icons.
var IconMap = map[string]string{

	// Premis Event Icons.
	// Note that we use only 7 or so event types,
	// so we don't define an icon for every type.
	constants.EventAccessAssignment:     `<i class="fas fa-user-shield"></i>`,
	constants.EventCreate:               `<i class="fas fa-plus-circle"></i>`,
	constants.EventDeletion:             `<i class="fas fa-minus-circle"></i>`,
	constants.EventDigestCalculation:    `<i class="fas fa-file-signature"></i>`,
	constants.EventFixityCheck:          `<i class="fas fa-fingerprint"></i>`,
	constants.EventIdentifierAssignment: `<i class="far fa-snowflake"></i>`,
	constants.EventIngestion:            `<i class="fas fa-file-import"></i>`,
	constants.EventReplication:          `<i class="fas fa-copy"></i>`,
}
