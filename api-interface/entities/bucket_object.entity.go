package entity

import (
	"encoding/json"
)

type BucketObject struct {
    Key          string          `json:"key"`
    LastModified string          `json:"lastModified"`
    ETag         string          `json:"etag"`
    Size         int64           `json:"size"`
    StorageClass StorageClassType `json:"storageClass"`
    Owner        Owner           `json:"owner"`
    Type         string          `json:"type"`
    URI          string          `json:"uri"`
    BucketName   string          `json:"bucketName"`
}

func (bo *BucketObject) GetKey() string {
    return bo.Key
}

func (bo *BucketObject) Serialize() ([]byte, error) {
    return json.Marshal(bo)
}

func (bo *BucketObject) Deserialize(data []byte) error {
    return json.Unmarshal(data, bo)
}