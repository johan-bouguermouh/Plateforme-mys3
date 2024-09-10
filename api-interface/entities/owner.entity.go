package entity

import (
	"encoding/json"
)

type Role struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Owner struct {
	UserKey     string `json:"UserKey"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
	ROLE        Role   `json:"role"`
	SecretKey   string `json:"secretKey"`
}

func (o *Owner) GetKey() string {
	return o.UserKey
}

func (o *Owner) Serialize() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Owner) Deserialize(data []byte) error {
    return json.Unmarshal(data, o)
}