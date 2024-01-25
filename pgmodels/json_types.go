package pgmodels

import (
	"time"
)

// This file contains structs and methods to convert pgmodel objects
// into subsets of the original structs that are safe for serialiation
// in the API. These are often used in the member API, where we don't
// want to expose as much info as we do in the admin API.
//
// For example, when serializing an object like a deletion request, we
// also serialize info about the users who created, approved, or
// cancelled the request. In this context, we do not want to serialize
// info about the user that we would typically serialize in a user management
// endpoint. This would include things like whether the user has registered
// for Authy two-factor authentication, where the last logged in from, etc.
//
// All we really need to serialize in such a case is the user's name, id, and
// email address.
//
// We don't want to write a custom JSON marshal method to do that, because
// Go only allows one JSON marshaller per type, and we need to serialize
// some of these objects two different ways. E.g. Serialize full user for
// the admin API, but serialize UserMin for most endpoints in the member
// API.

// UserMin contains minimum required information to identify a user.
// This is used for serialization in parts of the member API when we
// want to show which user is attached to a record (such as a deletion
// request) but we do not want to expose unnecessary information about
// that user.
type UserMin struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToMin returns a UserMin struct containing only enough info to identify
// a user. This subset of info is safe to expose through member API
// endpoints.
func (user *User) ToMin() *UserMin {
	return &UserMin{user.ID, user.Name, user.Email}
}

// DeletionRequestMin includes the subset of DeletionRequest info that
// we want to expose through the member API.
type DeletionRequestMin struct {
	ID                  int64                 `json:"id"`
	InstitutionID       int64                 `json:"institution_id"`
	RequestedAt         time.Time             `json:"requested_at"`
	ConfirmedAt         time.Time             `json:"confirmed_at"`
	CancelledAt         time.Time             `json:"cancelled_at"`
	RequestedBy         *UserMin              `json:"requested_by"`
	ConfirmedBy         *UserMin              `json:"confirmed_by"`
	CancelledBy         *UserMin              `json:"cancelled_by"`
	GenericFiles        []*GenericFile        `json:"generic_files"`
	IntellectualObjects []*IntellectualObject `json:"intellectual_objects"`
}

// ToMin returns DeletionRequestMin object suitable for serialization
// in the member API.
func (r *DeletionRequest) ToMin() *DeletionRequestMin {
	var confirmedBy *UserMin
	if r.ConfirmedBy != nil {
		confirmedBy = r.ConfirmedBy.ToMin()
	}
	var cancelledBy *UserMin
	if r.CancelledBy != nil {
		cancelledBy = r.CancelledBy.ToMin()
	}
	return &DeletionRequestMin{
		ID:                  r.ID,
		InstitutionID:       r.InstitutionID,
		RequestedAt:         r.RequestedAt,
		ConfirmedAt:         r.ConfirmedAt,
		CancelledAt:         r.CancelledAt,
		RequestedBy:         r.RequestedBy.ToMin(),
		ConfirmedBy:         confirmedBy,
		CancelledBy:         cancelledBy,
		GenericFiles:        r.GenericFiles,
		IntellectualObjects: r.IntellectualObjects,
	}
}
