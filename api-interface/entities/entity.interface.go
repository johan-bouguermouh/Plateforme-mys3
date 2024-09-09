package entity

// Interface destiné à être implémenté par les entités mangées par le système ORM et le QueryBuilder
type Entity interface {
	GetKey() string
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}