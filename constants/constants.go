package constants

const (
	AccessConsortia            = "consortia"
	AccessInstitution          = "institution"
	AccessRestricted           = "restricted"
	ActionApproveDelete        = "ApproveDelete"
	ActionCreate               = "Create"
	ActionDelete               = "Delete"
	ActionFinishBulkDelete     = "FinishBulkDelete"
	ActionFixityCheck          = "Fixity Check"
	ActionGlacierRestore       = "Glacier Restore"
	ActionIngest               = "Ingest"
	ActionRead                 = "Read"
	ActionRequestDelete        = "RequestDelete"
	ActionRestore              = "Restore"
	ActionUpdate               = "Update"
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
	TopicFileRestore           = "restore_file"
	TopicFixity                = "fixity_check"
	TopicGlacierRestore        = "restore_glacier"
	TopicObjectRestore         = "restore_object"
)

var UserActions = []string{
	ActionCreate,
	ActionRead,
	ActionUpdate,
	ActionDelete,
	ActionRequestDelete,
	ActionApproveDelete,
	ActionFinishBulkDelete,
	ActionRestore,
}

var WorkItemActions = []string{
	ActionDelete,
	ActionGlacierRestore,
	ActionIngest,
	ActionRestore,
}

var AccessSettings = []string{
	AccessConsortia,
	AccessInstitution,
	AccessRestricted,
}

var DigestAlgs = []string{
	AlgMd5,
	AlgSha1,
	AlgSha256,
	AlgSha512,
}

var InstTypes = []string{
	InstTypeMember,
	InstTypeSubscriber,
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
