package pgmodels_test

import (
	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/pgmodels"
)

// getTestObject returns an IntellectualObject with valid settings
// that can be altered per-test.
func GetTestObject() *pgmodels.IntellectualObject {
	return &pgmodels.IntellectualObject{
		Title:                     "TestObject999",
		Description:               "Obj Created by Test",
		Identifier:                "test.edu/obj1",
		AltIdentifier:             "Yadda-Yadda-Yo",
		Access:                    constants.AccessInstitution,
		State:                     constants.StateActive,
		BagName:                   "TestObject999.tar",
		ETag:                      "12345678-9",
		InstitutionID:             4,
		StorageOption:             constants.StorageOptionStandard,
		BagItProfileIdentifier:    "https://example.com/profile.json",
		SourceOrganization:        "Willy Wonka's Chocolate Factory",
		BagGroupIdentifier:        "group-999",
		InternalSenderIdentifier:  "yadda-999",
		InternalSenderDescription: "Created by intel obj test",
	}
}
