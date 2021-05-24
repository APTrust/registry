package constants

import "github.com/APTrust/registry/common"

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
	AlertRestorationCompleted  = "Restoration Completed"
	AlertStalledItems          = "Stalled Work Items"
	AlgMd5                     = "md5"
	AlgSha1                    = "sha1"
	AlgSha256                  = "sha256"
	AlgSha512                  = "sha512"
	DefaultProfileIdentifier   = "https://raw.githubusercontent.com/APTrust/preservation-services/master/profiles/aptrust-v2.2.json"
	EventAccessAssignment      = "access assignment"
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
	RoleInstAdmin              = "institutional_admin"
	RoleInstUser               = "institutional_user"
	RoleNone                   = "none"
	RoleSysAdmin               = "admin"
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
	StorageOptionWasabiVA      = "Wasabi-VA"
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
	AlertStalledItems,
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
	ActionRestoreObject,
	ActionRestoreFile,
	ActionGlacierRestore,
	ActionIngest,
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
		err = common.ErrInvalidRequeue
	}
	return topic, err
}
