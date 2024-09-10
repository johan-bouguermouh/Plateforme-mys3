package entity

import (
	"encoding/json"
)

type Bucket struct {
    Name         string          `json:"name"`
    CreationDate string          `json:"creationDate"`
    Owner        Owner           `json:"owner"`
    URI          string          `json:"uri"`
    Type         BucketType      `json:"type"`
    StorageClass StorageClassType `json:"storageClass"`
    Versioning   VersioningStatus `json:"versioning"`
    ObjectCount  int64           `json:"objectCount"`
    Size         int64           `json:"size"`
    LastModified string          `json:"lastModified"`
}

func (b *Bucket) GetKey() string {
    return b.Name
}

func (b *Bucket) Serialize() ([]byte, error) {
    return json.Marshal(b)
}

func (b *Bucket) Deserialize(data []byte) error {
    return json.Unmarshal(data, b)
}