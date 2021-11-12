package admin_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IntellectualObjectCreate creates a new object record.
//
// POST /admin-api/v3/objects/create/:institution_id
func IntellectualObjectCreate(c *gin.Context) {
	// Ensure the inst id in the JSON matches what's in the URL
	// Create the object.
	// Return the full object record.
	c.JSON(http.StatusOK, nil)
}

// IntellectualObjectUpdate updates an existing intellectual
// object record.
//
// PUT /admin-api/v3/objects/update/:id
func IntellectualObjectUpdate(c *gin.Context) {
	// Ensure the inst id in the JSON matches what's in the URL
	// Update the object, ensuring:
	//  - institution id can't change
	//  - storage option can't change
	// Return the full object record.
	c.JSON(http.StatusOK, nil)
}

// IntellectualObjectDelete marks an object record as deleted.
//
// DELETE /admin-api/v3/objects/delete/:id
func IntellectualObjectDelete(c *gin.Context) {
	// We should probably not allow the object to be deleted
	// unless all of its files have been deleted. Double check
	// the business logic in Pharos.
	//
	// Object deletion changes the state from "A" to "D".
	//
	// We should also ensure a Premis Event exists or is created
	// the documents who deleted this and when.
	//
	// Check the Pharos logic on that too. It may be the Go
	// worker's responsibility to ensure this, or it may be
	// registry's responsibility.
	//
	// Return the full object record.
	c.JSON(http.StatusOK, nil)
}
