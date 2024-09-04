package models

// User model
type User struct {
	Name string `json:"name"`
}

type Role struct {
	/** Role name */
	Name string `json:"name"`
	/** Role ID */
	ID string `json:"id"`
	/** Role type */
	Type string `json:"type"`
}

/** An object type for find Owner of file in bucket S3 */
type Owner struct {
	/** User Key*/
	UserKey string `json:"UserKey"`
	/** Name of owner */
	DisplayName string `json:"displayName"`
	/** Type of owner */
	Type string `json:"type"`
	/** URI of owner */
	URI string `json:"uri"`
	/** ROLE of owner */
	ROLE Role `json:"role"`
	/** Secret Key of owner */
	SecretKey string `json:"secretKey"`
}
type StorageClassType string

const (
    STANDARD              StorageClassType = "STANDARD"
    REDUCED_REDUNDANCY    StorageClassType = "REDUCED_REDUNDANCY"
    STANDARD_IA           StorageClassType = "STANDARD_IA"
    ONEZONE_IA            StorageClassType = "ONEZONE_IA"
    INTELLIGENT_TIERING   StorageClassType = "INTELLIGENT_TIERING"
    GLACIER               StorageClassType = "GLACIER"
    DEEP_ARCHIVE          StorageClassType = "DEEP_ARCHIVE"
)

type BucketType string
const (
    PUBLIC  BucketType = "PUBLIC"
    PRIVATE BucketType = "PRIVATE"
)

type VersioningStatus string
const (
    VersioningEnabled  VersioningStatus = "Enabled"
    VersioningSuspended VersioningStatus = "Suspended"
    VersioningDisabled  VersioningStatus = "Disabled"
)

/** An object type for find Bucket in S3 */
type Bucket struct {
    /** Name of bucket */
    Name string `json:"name"`
    /** Date of creation */
    CreationDate string `json:"creationDate"`
    /** Owner of bucket */
    Owner Owner `json:"owner"`
    /** URI of bucket */
    URI string `json:"uri"`
    /** Type of bucket */
    Type BucketType `json:"type"`
    /**  Storage class, use for storage type : STANDARD, REDUCED_REDUNDANCY, STANDARD_IA, ONEZONE_IA, INTELLIGENT_TIERING, GLACIER, DEEP_ARCHIVE */
    StorageClass StorageClassType `json:"storageClass"`
    /** Versioning, use for versioning file */
    Versioning VersioningStatus `json:"versioning"`
    /** Encrytion, use for encryption file */
    ObjectCount int64 `json:"objectCount"`
    /** Size of bucket */
    Size int64 `json:"size"`
    /** Last modified date, use for check file integrity */
    LastModified string `json:"lastModified"`
}

/** An object type for find File in bucket S3 */
type BucketObject struct {
	/** Key name for find file */
	Key string `json:"key"`
	/** Last modified date, use for check file integrity */
	LastModified string `json:"lastModified"`
	/** ETag, use for cashing file and identifier for a specific version of a resource" */
	ETag string `json:"eTag"`
	/** Size of file */
	Size int64 `json:"size"`
	/** Storage class, use for storage type : STANDARD, REDUCED_REDUNDANCY, STANDARD_IA, ONEZONE_IA, INTELLIGENT_TIERING, GLACIER, DEEP_ARCHIVE */
	StorageClass StorageClassType `json:"storageClass"`
	/** Owner of file */
	Owner Owner `json:"owner"`
	/** Type of file, use for check file type : OBJECT, DIRECTORY */
	Type string `json:"type"`
	/** URI of file */
	URI string `json:"uri"`
	/** Bucket name of file */
	BucketName string `json:"bucketName"`
}