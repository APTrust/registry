package constants

import (
	"errors"
)

const (
	AccessConsortia            = "consortia"
	AccessInstitution          = "institution"
	AccessRestricted           = "restricted"
	ActionApproveDelete        = "ApproveDelete"
	ActionCreate               = "Create"
	ActionDelete               = "Delete"
	ActionFinishBulkDelete     = "FinishBulkDelete"
	ActionFixityCheck          = "Fixity Check"
	ActionRestoreFile          = "Restore File"
	ActionGlacierRestore       = "Glacier Restore"
	ActionIngest               = "Ingest"
	ActionRead                 = "Read"
	ActionRequestDelete        = "RequestDelete"
	ActionRestoreObject        = "Restore Object"
	ActionUpdate               = "Update"
	AlertDeletionCancelled     = "Deletion Cancelled"
	AlertDeletionCompleted     = "Deletion Completed"
	AlertDeletionConfirmed     = "Deletion Confirmed"
	AlertDeletionRequested     = "Deletion Requested"
	AlertFailedFixity          = "Failed Fixity Check"
	AlertPasswordChanged       = "Password Changed"
	AlertPasswordReset         = "Password Reset"
	AlertRestorationCompleted  = "Restoration Completed"
	AlertStalledItems          = "Stalled Work Items"
	AlertWelcome               = "Welcome New User"
	AlgMd5                     = "md5"
	AlgSha1                    = "sha1"
	AlgSha256                  = "sha256"
	AlgSha512                  = "sha512"
	APIUserHeader              = "X-Pharos-API-User"
	APIKeyHeader               = "X-Pharos-API-Key"
	APIPrefixAdmin             = "/admin-api/"
	APIPrefixMember            = "/member-api/"
	APTrustOpsEmail            = "ops@aptrust.org"
	BTRProfileIdentifier       = "https://github.com/dpscollaborative/btr_bagit_profile/releases/download/1.0/btr-bagit-profile.json"
	CSRFCookieName             = "csrf_token"
	CSRFHeaderName             = "X-CSRF-Token"
	CSRFTokenName              = "csrf_token"
	DefaultProfileIdentifier   = "https://raw.githubusercontent.com/APTrust/preservation-services/master/profiles/aptrust-v2.2.json"
	EmailServiceSES            = "SES"
	EmailServiceSMTP           = "SMTP"
	EventAccessAssignment      = "access assignment"
	EventAgentMinioV4          = 1
	EventAgentStringMinioV4    = "https://github.com/minio/minio-go v4"
	EventAgentMinioV5          = 2
	EventAgentStringMinioV5    = "https://github.com/minio/minio-go v5"
	EventAgentMinioV6          = 3
	EventAgentStringMinioV6    = "https://github.com/minio/minio-go v6"
	EventAgentMinioV7          = 4
	EventAgentStringMinioV7    = "https://github.com/minio/minio-go v7"
	EventAgentPreserv          = 5
	EventAgentStringPreserv    = "https://github.com/APTrust/preservation-services"
	EventAgentPreservAlt       = 6
	EventAgentStringPreservAlt = "APTrust preservation services"
	EventAgentUUID             = 7
	EventAgentStringUUID       = "http://github.com/google/uuid"
	EventAgentSHA256           = 8
	EventAgentStringSHA256     = "http://golang.org/pkg/crypto/sha256/"
	EventAgentMD5              = 9
	EventAgentStringMD5        = "http://golang.org/pkg/crypto/md5/"
	EventAgentTest             = 10
	EventAgentStringTest       = "Registry Unit Test"
	EventAgentTestAlt          = 11
	EventAgentStringTestAlt    = "Maxwell Smart"
	EventAgentFixture          = 12
	EventAgentStringFixture    = "https://github.com/APTrust/exchange"
	EventAgentUUIDPast         = 13
	EventAgentStringUUIDPast   = "http://github.com/satori/go.uuid"
	EventAgentNu7Hatch         = 14
	EventAgentStringNu7Hatch   = "http://github.com/nu7hatch/gouuid"
	EventAgentBagman           = 15
	EventAgentStringBagman     = "https://github.com/APTrust/bagman"
	EventAgentAwsSdk           = 16
	EventAgentStringAwsSdk     = "https://github.com/aws/aws-sdk-go"
	EventAgentGoamz            = 17
	EventAgentStringGoamz      = "https://github.com/crowdmob/goamz"
	EventAgentMarcel           = 18
	EventAgentStringMarcel     = "https://github.com/marcel/aws-s3/tree/master"
	EventAgentUuidHttps        = 19
	EventAgentStringUuidHttps  = "https://github.com/satori/go.uuid"
	EventAgentLaunchpad        = 20
	EventAgentStringLaunchpad  = "https://launchpad.net/goamz"
	EventAgentAudit            = 21
	EventAgentStringAudit      = "https://github.com/APTrust/auditing/blob/1.0/cleanup_001.py"
	EventCapture               = "capture"
	EventCompression           = "compression"
	EventCreation              = "creation"
	EventDeaccession           = "deaccession"
	EventDecompression         = "decompression"
	EventDecryption            = "decryption"
	EventDeletion              = "deletion"
	EventDigestCalculation     = "message digest calculation"
	EventFixityCheck           = "fixity check"
	EventIdentifierAssignment  = "identifier assignment"
	EventIngestion             = "ingestion"
	EventMigration             = "migration"
	EventNormalization         = "normalization"
	EventObjectPreserv         = 1
	EventObjectStringPreserv   = "APTrust preservation services"
	EventObjectMinio           = 2
	EventObjectStringMinio     = "Minio S3 client"
	EventObjectMinioAlt        = 3
	EventObjectStringMinioAlt  = "Minio S3 library"
	EventObjectPreMinio        = 4
	EventObjectStringPreMinio  = "preservation-services + Minio S3 client"
	EventObjectUUIDMinio       = 5
	EventObjectStringUUIDMinio = "Go uuid library + Minio S3 library"
	EventObjectSHA256          = 6
	EventObjectStringSHA256    = "Go language crypto/sha256"
	EventObjectMD5             = 7
	EventObjectStringMD5       = "Go language crypto/md5"
	EventObjectTest            = 8
	EventObjectStringTest      = "scissors"
	EventObjectExchange        = 9
	EventObjectStringExchange  = "APTrust exchange/ingest processor"
	EventObjectTestAlt         = 10
	EventObjectStringTestAlt   = "Fake event object"
	EventObjectFixS3           = 11
	EventObjectStringFixS3     = "APTrust Go Exchange + Amazon S3 client"
	EventObjectFixSHA          = 12
	EventObjectStringFixSHA    = "SHA-256 thingy"
	EventObjectFixExch         = 13
	EventObjectStringFixExch   = "Exchange ingest code"
	EventObjectFixDelete       = 14
	EventObjectStringFixDelete = "Deleterbot code"
	EventObjectBagman          = 15
	EventObjectStringBagman    = "APTrust bagman"
	EventObjectBagProc         = 16
	EventObjectStringBagProc   = "APTrust bag processor"
	EventObjectExchOld         = 17
	EventObjectStringExchOld   = "APTrust exchange"
	EventObjectExDelete        = 18
	EventObjectStringExDelete  = "APTrust Exchange apt_delete service"
	EventObjectExIngest        = 19
	EventObjectStringExIngest  = "APTrust Exchange ingest services"
	EventObjectExUUID          = 20
	EventObjectStringExUUID    = "APTrust exchange using Satori go.uuid"
	EventObjectAwsClient       = 21
	EventObjectStringAwsClient = "AWS Go SDK S3 client"
	EventObjectAwsLib          = 22
	EventObjectStringAwsLib    = "AWS Go SDK S3 Library"
	EventObjectBagmanGo        = 23
	EventObjectStringBagmanGo  = "bagman + goamz s3 client"
	EventObjectExAws           = 24
	EventObjectStringExAws     = "exchange + AWS Go SDK S3 client"
	EventObjectExGoamz         = 25
	EventObjectStringExGoamz   = "exchange + goamz S3 client"
	EventObjectGoamz           = 26
	EventObjectStringGoamz     = "goamz S3 client"
	EventObjectGoamzAlt        = 27
	EventObjectStringGoamzAlt  = "Goamz S3 Client"
	EventObjectCryptoh         = 28
	EventObjectStringCryptoh   = "Go language cryptohash"
	EventObjectMD5Past         = 29
	EventObjectStringMD5Past   = "Go crypto/md5"
	EventObjectDPN             = 30
	EventObjectStringDPN       = "Go uuid library + APTrust DPN services"
	EventObjectUuidAws         = 31
	EventObjectStringUuidAws   = "Go uuid library + AWS Go SDK S3 library"
	EventObjectUuidGoamz       = 32
	EventObjectStringUuidGoamz = "Go uuid library + goamz S3 library"
	EventObjectRuby            = 33
	EventObjectStringRuby      = "Ruby aws-s3 gem"
	EventObjectAudit           = 34
	EventObjectStringAudit     = "APTrust audit and cleanup scripts for audit_001"
	EventReplication           = "replication"
	EventSignatureValidation   = "digital signature validation"
	EventValidation            = "validation"
	EventVirusCheck            = "virus check"
	IngestPreFetch             = "ingest01_prefetch"
	IngestValidation           = "ingest02_bag_validation"
	IngestReingestCheck        = "ingest03_reingest_check"
	IngestStaging              = "ingest04_staging"
	IngestFormatIdentification = "ingest05_format_identification"
	IngestStorage              = "ingest06_storage"
	IngestStorageValidation    = "ingest07_storage_validation"
	IngestRecord               = "ingest08_record"
	IngestCleanup              = "ingest09_cleanup"
	InstTypeMember             = "MemberInstitution"
	InstTypeSubscriber         = "SubscriptionInstitution"
	MetaFixityAlertsLastRun    = "fixity alerts last run"
	MetaSpotTestsRunning       = "spot restore is running"
	MetaSpotTestsLastRun       = "spot restore last run"
	OutcomeFailure             = "Failure"
	OutcomeSuccess             = "Success"
	RoleInstAdmin              = "institutional_admin"
	RoleInstUser               = "institutional_user"
	RoleNone                   = "none"
	RoleSysAdmin               = "admin"
	SecondFactorTOTP           = "Authenticator App"
	SecondFactorBackupCode     = "Backup Code"
	SecondFactorSMS            = "SMS"
	StageAvailableInS3         = "Available in S3"
	StageCleanup               = "Cleanup"
	StageCopyToStaging         = "Copy To Staging"
	StageFetch                 = "Fetch"
	StageFormatIdentification  = "Format Identification"
	StagePackage               = "Package"
	StageReceive               = "Receive"
	StageRecord                = "Record"
	StageReingestCheck         = "Reingest Check"
	StageRequested             = "Requested"
	StageResolve               = "Resolve"
	StageRestoring             = "Restoring"
	StageStorageValidation     = "Storage Validation"
	StageStore                 = "Store"
	StageUnpack                = "Unpack"
	StageValidate              = "Validate"
	StateActive                = "A"
	StateDeleted               = "D"
	StatusCancelled            = "Cancelled"
	StatusFailed               = "Failed"
	StatusPending              = "Pending"
	StatusStarted              = "Started"
	StatusSuccess              = "Success"
	StatusSuspended            = "Suspended"
	StorageOptionGlacierDeepOH = "Glacier-Deep-OH"
	StorageOptionGlacierDeepOR = "Glacier-Deep-OR"
	StorageOptionGlacierDeepVA = "Glacier-Deep-VA"
	StorageOptionGlacierOH     = "Glacier-OH"
	StorageOptionGlacierOR     = "Glacier-OR"
	StorageOptionGlacierVA     = "Glacier-VA"
	StorageOptionStandard      = "Standard"
	StorageOptionWasabiOR      = "Wasabi-OR"
	StorageOptionWasabiTX      = "Wasabi-TX"
	StorageOptionWasabiVA      = "Wasabi-VA"
	SystemUser                 = "system@aptrust.org"
	TopicDelete                = "delete_item"
	TopicE2EDelete             = "e2e_deletion_post_test"
	TopicE2EFixity             = "e2e_fixity_post_test"
	TopicE2EIngest             = "e2e_ingest_post_test"
	TopicE2EReingest           = "e2e_reingest_post_test"
	TopicE2ERestore            = "e2e_restoration_post_test"
	TopicFileRestore           = "restore_file"
	TopicFixity                = "fixity_check"
	TopicGlacierRestore        = "restore_glacier"
	TopicObjectRestore         = "restore_object"
	TOTPSecretIssuer           = "APTrust"
	TwoFactorTOTP              = "totp"
	TwoFactorNone              = "none"
	TwoFactorSMS               = "sms"
)

var AccessSettings = []string{
	AccessConsortia,
	AccessInstitution,
	AccessRestricted,
}

var AlertTypes = []string{
	AlertDeletionCancelled,
	AlertDeletionCompleted,
	AlertDeletionConfirmed,
	AlertDeletionRequested,
	AlertFailedFixity,
	AlertRestorationCompleted,
	AlertPasswordChanged,
	AlertPasswordReset,
	AlertStalledItems,
	AlertWelcome,
}

var APIPrefixes = []string{
	APIPrefixAdmin,
	APIPrefixMember,
}

var CompletedStatusValues = []string{
	StatusCancelled,
	StatusFailed,
	StatusSuccess,
}

var DigestAlgs = []string{
	AlgMd5,
	AlgSha1,
	AlgSha256,
	AlgSha512,
}

var EventOutcomes = []string{
	OutcomeFailure,
	OutcomeSuccess,
}

var EventTypes = []string{
	EventAccessAssignment,
	EventCreation,
	EventDeletion,
	EventDigestCalculation,
	EventFixityCheck,
	EventIdentifierAssignment,
	EventIngestion,
	EventReplication,
	EventValidation,
}

var GlacierOnlyOptions = []string{
	StorageOptionGlacierDeepOH,
	StorageOptionGlacierDeepOR,
	StorageOptionGlacierDeepVA,
	StorageOptionGlacierOH,
	StorageOptionGlacierOR,
	StorageOptionGlacierVA,
}

var IncompleteStatusValues = []string{
	StatusPending,
	StatusStarted,
}

var IngestStagesInOrder = []string{
	StageReceive,
	StageValidate,
	StageReingestCheck,
	StageCopyToStaging,
	StageFormatIdentification,
	StageStore,
	StageStorageValidation,
	StageRecord,
	StageCleanup,
}

var InstTypes = []string{
	InstTypeMember,
	InstTypeSubscriber,
}

var Roles = []string{
	RoleInstAdmin,
	RoleInstUser,
	RoleNone,
	RoleSysAdmin,
}

var SecondFactorTypes = []string{
	SecondFactorBackupCode,
	SecondFactorSMS,
}

var Stages = []string{
	StageAvailableInS3,
	StageCleanup,
	StageCopyToStaging,
	StageFormatIdentification,
	StageFetch,
	StagePackage,
	StageReceive,
	StageRecord,
	StageReingestCheck,
	StageRequested,
	StageResolve,
	StageRestoring,
	StageStorageValidation,
	StageStore,
	StageUnpack,
	StageValidate,
}

var States = []string{
	StateActive,
	StateDeleted,
}

var Statuses = []string{
	StatusCancelled,
	StatusFailed,
	StatusPending,
	StatusStarted,
	StatusSuccess,
	StatusSuspended,
}

var StorageOptions = []string{
	StorageOptionGlacierDeepOH,
	StorageOptionGlacierDeepOR,
	StorageOptionGlacierDeepVA,
	StorageOptionGlacierOH,
	StorageOptionGlacierOR,
	StorageOptionGlacierVA,
	StorageOptionStandard,
	StorageOptionWasabiOR,
	StorageOptionWasabiTX,
	StorageOptionWasabiVA,
}

var UserActions = []string{
	ActionCreate,
	ActionRead,
	ActionUpdate,
	ActionDelete,
	ActionRequestDelete,
	ActionApproveDelete,
	ActionFinishBulkDelete,
	ActionRestoreObject,
	ActionRestoreFile,
}

var WorkItemActions = []string{
	ActionDelete,
	ActionGlacierRestore,
	ActionIngest,
	ActionRestoreFile,
	ActionRestoreObject,
}

var NonIngestTopics = []string{
	TopicDelete,
	TopicFileRestore,
	TopicFixity,
	TopicGlacierRestore,
	TopicObjectRestore,
}

// NSQIngestTopicFor maps ingest stage names to NSQ topics.
var NSQIngestTopicFor = map[string]string{
	StageReceive:              IngestPreFetch,
	StageValidate:             IngestValidation,
	StageReingestCheck:        IngestReingestCheck,
	StageCopyToStaging:        IngestStaging,
	StageFormatIdentification: IngestFormatIdentification,
	StageStore:                IngestStorage,
	StageStorageValidation:    IngestStorageValidation,
	StageRecord:               IngestRecord,
	StageCleanup:              IngestCleanup,
}

// All other common errors are defined in common. We had to move this
// one into constants to prevent an illegal import cycle.

// ErrInvalidRequeue occurs when someone attempts to requeue an item to the
// wrong stage, or to a stage for which no NSQ topic exists.
var ErrInvalidRequeue = errors.New("item cannot be requeued to the specified stage")

func TopicFor(action, stage string) (string, error) {
	var err error
	topic := ""
	switch action {
	case ActionDelete:
		topic = TopicDelete
	case ActionRestoreFile:
		topic = TopicFileRestore
	case ActionGlacierRestore:
		topic = TopicGlacierRestore
	case ActionRestoreObject:
		topic = TopicObjectRestore
	case ActionIngest:
		topic = NSQIngestTopicFor[stage]
	default:
		topic = ""
	}
	if topic == "" {
		err = ErrInvalidRequeue
	}
	return topic, err
}
