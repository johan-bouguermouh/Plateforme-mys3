package entity

type StorageClassType string

const (
	STANDARD            StorageClassType = "STANDARD"
	REDUCED_REDUNDANCY  StorageClassType = "REDUCED_REDUNDANCY"
	STANDARD_IA         StorageClassType = "STANDARD_IA"
	ONEZONE_IA          StorageClassType = "ONEZONE_IA"
	INTELLIGENT_TIERING StorageClassType = "INTELLIGENT_TIERING"
	GLACIER             StorageClassType = "GLACIER"
	DEEP_ARCHIVE        StorageClassType = "DEEP_ARCHIVE"
)

type BucketType string

const (
	PUBLIC  BucketType = "PUBLIC"
	PRIVATE BucketType = "PRIVATE"
)

type VersioningStatus string

const (
	VersioningEnabled   VersioningStatus = "Enabled"
	VersioningSuspended VersioningStatus = "Suspended"
	VersioningDisabled  VersioningStatus = "Disabled"
)