package pgmodels

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

// --------------------------------------------------------------------
// TODO:
//
// ✓ Explain the magic number 4 & set as constant.
// x  Rename these methods so it's clear they are factory methods
//    to be used for testing only. (Maybe move them to a different
//    package.)
// --------------------------------------------------------------------

var testInstitutionID = int64(4)
var testInstAdminID = int64(8)
var testInstUserID = int64(9)

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
		GenericFiles:              make([]*GenericFile, 0),
		PremisEvents:              make([]*PremisEvent, 0),
	}
}

func CreateObjectWithRelations() (*IntellectualObject, error) {
	obj := RandomObject()
	err := obj.Save()
	if err != nil {
		goto ERR
	}
	for i := 0; i < 5; i++ {
		gf := RandomGenericFile(obj.ID, obj.Identifier)
		err = gf.Save()
		if err != nil {
			goto ERR
		}
		gf.PremisEvents = append(gf.PremisEvents, RandomPremisEvent(constants.EventIngestion))
		gf.PremisEvents = append(gf.PremisEvents, RandomPremisEvent(constants.EventIdentifierAssignment))
		for _, event := range gf.PremisEvents {
			event.IntellectualObjectID = gf.IntellectualObjectID
			event.GenericFileID = gf.ID
			err = event.Save()
			if err != nil {
				goto ERR
			}
		}
		gf.Checksums = append(gf.Checksums, RandomChecksum(constants.AlgMd5))
		gf.Checksums = append(gf.Checksums, RandomChecksum(constants.AlgSha256))
		for _, checksum := range gf.Checksums {
			checksum.GenericFileID = gf.ID
			err = checksum.Save()
			if err != nil {
				goto ERR
			}
		}

		gf.StorageRecords = append(gf.StorageRecords, RandomStorageRecord())
		gf.StorageRecords = append(gf.StorageRecords, RandomStorageRecord())
		for _, sr := range gf.StorageRecords {
			sr.GenericFileID = gf.ID
			err = sr.Save()
			if err != nil {
				goto ERR
			}
		}
		obj.GenericFiles = append(obj.GenericFiles, gf)
	}

	obj.PremisEvents = append(obj.PremisEvents, RandomPremisEvent(constants.EventIngestion))
	obj.PremisEvents = append(obj.PremisEvents, RandomPremisEvent(constants.EventIdentifierAssignment))
	for _, event := range obj.PremisEvents {
		event.IntellectualObjectID = obj.ID
		err = event.Save()
		if err != nil {
			goto ERR
		}
	}

	return obj, nil
ERR:
	return nil, err
}

func RandomWorkItem(name, action string, objID, gfID int64) *WorkItem {
	now := time.Now().UTC()
	return &WorkItem{
		TimestampModel: TimestampModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:                 name,
		ETag:                 ETag(),
		InstitutionID:        4,
		IntellectualObjectID: objID,
		GenericFileID:        gfID,
		Bucket:               "blah.receiving.blah",
		User:                 "someone@example.com",
		Note:                 "This item was created by the factory",
		Action:               action,
		Stage:                constants.StageRequested,
		Status:               constants.StatusPending,
		Outcome:              "Outcome? WTF?",
		BagDate:              now.AddDate(0, -4, 0),
		DateProcessed:        now,
		Retry:                true,
		QueuedAt:             now,
		Size:                 rand.Int63(),
	}
}

// RandomGenericFile returns a random generic file with the specified
// obj identifier prefix. State will be active.
func RandomGenericFile(objID int64, objIdentifier string) *GenericFile {
	now := time.Now().UTC()
	return &GenericFile{
		TimestampModel: TimestampModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		FileFormat:           gofakeit.AnimalType(),
		Size:                 rand.Int63(),
		Identifier:           FileIdentifier(objIdentifier),
		InstitutionID:        4,
		IntellectualObjectID: objID,
		State:                constants.StateActive,
		LastFixityCheck:      now.AddDate(0, -4, 0),
		StorageOption:        constants.StorageOptionStandard,
		UUID:                 uuid.NewString(),
		PremisEvents:         make([]*PremisEvent, 0),
		Checksums:            make([]*Checksum, 0),
		StorageRecords:       make([]*StorageRecord, 0),
	}
}

// RandomPremisEvent returns a random premis event of the specified
// type. Caller should set GenericFileID and IntellectualObjectID.
func RandomPremisEvent(eventType string) *PremisEvent {
	now := time.Now().UTC()
	return &PremisEvent{
		Agent:              gofakeit.FarmAnimal(),
		DateTime:           now,
		Detail:             gofakeit.Sentence(4),
		EventType:          eventType,
		Identifier:         uuid.NewString(),
		InstitutionID:      4,
		Object:             gofakeit.Sentence(4),
		Outcome:            gofakeit.Sentence(5),
		OutcomeDetail:      gofakeit.BeerName(),
		OutcomeInformation: gofakeit.AppAuthor(),
		TimestampModel: TimestampModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// RandomChecksum returns a random checksum with the specified
// algorithm. Caller should set GenericFileID.
func RandomChecksum(alg string) *Checksum {
	now := time.Now().UTC()
	return &Checksum{
		TimestampModel: TimestampModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Algorithm: alg,
		DateTime:  now,
		Digest:    ETag(),
	}
}

// RandomStorageRecord() returns a random storage record.
// Caller should set GenericFileID.
func RandomStorageRecord() *StorageRecord {
	return &StorageRecord{
		URL: gofakeit.URL(),
	}
}

// RandFileBatch returns a slice of 20 GenericFiles, each with
// 4 events, checksums, and storage records. Note that this also
// creates the parent object in the database. You may want to
// reset the DB after calling this.
func RandomFileBatch() (*IntellectualObject, []*GenericFile, error) {
	obj := RandomObject()
	err := obj.Save()
	if err != nil {
		return nil, nil, err
	}

	files := make([]*GenericFile, 20)
	for i := 0; i < 20; i++ {
		gf := RandomGenericFile(obj.ID, obj.Identifier)
		for j := 0; j < 4; j++ {
			event := RandomPremisEvent(constants.EventIngestion)
			event.Outcome = constants.OutcomeSuccess
			gf.PremisEvents = append(gf.PremisEvents, event)

			checksum := RandomChecksum(constants.AlgSha1)
			gf.Checksums = append(gf.Checksums, checksum)

			sr := RandomStorageRecord()
			gf.StorageRecords = append(gf.StorageRecords, sr)
		}
		files[i] = gf
	}
	return obj, files, nil
}

func CreateDeletionRequest(objects []*IntellectualObject, files []*GenericFile) (*DeletionRequest, error) {
	now := time.Now().UTC()
	req, err := NewDeletionRequest()
	if err != nil {
		return nil, err
	}
	req.InstitutionID = 4
	req.RequestedAt = now.Add(-1 * time.Hour)
	req.RequestedByID = testInstUserID
	req.GenericFiles = files
	req.IntellectualObjects = objects

	err = req.Save()
	return req, err
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
