package pgmodels

import (
	"fmt"
	"strings"

	"github.com/APTrust/registry/constants"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

// GetTestObject returns an IntellectualObject with valid settings
// that can be altered per-test.
func GetTestObject() *IntellectualObject {
	return &IntellectualObject{
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

func RandomObject() *IntellectualObject {
	return &IntellectualObject{
		Title:                     Title(),
		Description:               gofakeit.HackerPhrase(),
		Identifier:                ObjIdentifier(),
		AltIdentifier:             gofakeit.FarmAnimal(),
		Access:                    constants.AccessInstitution,
		State:                     constants.StateActive,
		BagName:                   BagName(),
		ETag:                      ETag(),
		InstitutionID:             4,
		StorageOption:             constants.StorageOptionStandard,
		BagItProfileIdentifier:    "https://example.com/profile.json",
		SourceOrganization:        "Test University",
		BagGroupIdentifier:        gofakeit.Noun(),
		InternalSenderIdentifier:  gofakeit.PetName(),
		InternalSenderDescription: gofakeit.HipsterSentence(12),
	}
}

func Title() string {
	return fmt.Sprintf("%s %s %s", gofakeit.HackerAdjective(), gofakeit.BuzzWord(), gofakeit.BS())
}

func BagName() string {
	return fmt.Sprintf("%s.tar", gofakeit.HackerAdjective())
}

func ObjIdentifier() string {
	return fmt.Sprintf("test.edu/%s", gofakeit.AppName())
}

func FileIdentifier(objIdentifier string) string {
	return fmt.Sprintf("%s/data/%s.%s", objIdentifier, gofakeit.Gamertag(), gofakeit.FileExtension())
}

func ETag() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}
