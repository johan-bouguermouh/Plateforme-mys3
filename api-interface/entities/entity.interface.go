package entity

type Entity interface {
	GetKey() string
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}