package entity

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// Initialisation des types d'entités
func init() {
    RegisterEntityType("Bucket", reflect.TypeOf(Bucket{}))
    RegisterEntityType("BucketObject", reflect.TypeOf(BucketObject{}))
    RegisterEntityType("Owner", reflect.TypeOf(Owner{}))
}

// EntityRegistry est une structure qui contient une carte des types d'entités enregistrés.
type EntityRegistry struct {
    mu    sync.RWMutex
    types map[string]reflect.Type
}

// FieldInfo contient des informations sur un champ d'une entité
type FieldInfo struct {
    Name string
    Type reflect.Type
    Tag  reflect.StructTag
}

// NewEntityRegistry crée une nouvelle instance de EntityRegistry.
func NewEntityRegistry() *EntityRegistry {
    return &EntityRegistry{
        types: make(map[string]reflect.Type),
    }
}

// Register enregistre un type d'entité dans la carte.
func (r *EntityRegistry) Register(name string, t reflect.Type) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.types[name] = t
}

// Get récupère un type d'entité à partir de la carte.
func (r *EntityRegistry) Get(name string) (reflect.Type, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    t, ok := r.types[name]
    return t, ok
}

// Global registry instance
var globalRegistry = NewEntityRegistry()

// RegisterEntityType enregistre un type d'entité dans la carte globale.
func RegisterEntityType(name string, t reflect.Type) {
    globalRegistry.Register(name, t)
}

// GetEntityType récupère un type d'entité à partir de la carte globale.
func GetEntityType(name string) (reflect.Type, bool) {
    return globalRegistry.Get(name)
}

// GetAllEntityTypes retourne tous les types d'entités enregistrés dans le registre global.
func GetAllEntityTypes() map[string]reflect.Type {
    globalRegistry.mu.RLock()
    defer globalRegistry.mu.RUnlock()
    return globalRegistry.types
}

// GetFieldsByEntityType retourne tous les attributs présents dans une entité
func GetFieldsByEntityType(name string) ([]FieldInfo, error) {
    structType, exists := GetEntityType(name)
    if !exists {
        return nil, fmt.Errorf("entity type %s not found", name)
    }

    var fields []FieldInfo
    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        fields = append(fields, FieldInfo{
            Name: field.Name,
            Type: field.Type,
            Tag:  field.Tag,
        })
    }
    return fields, nil
}

// AttributeExists vérifie si un attribut existe dans une entité
func AttributeExists(entityName string, attributeName string) (bool, error) {
    fields, err := GetFieldsByEntityType(entityName)
    if err != nil {
        return false, err
    }

    for _, field := range fields {
        if field.Name == attributeName {
            return true, nil
        }
    }
    return false, nil
}

// TypeOfAttributeName retourne le type d'un attribut dans une entité
func TypeOfAttributeName(entityName string, attributeName string) (reflect.Type, error) {
    fields, err := GetFieldsByEntityType(entityName)
    if err != nil {
        return nil, err
    }

    for _, field := range fields {
        if field.Name == attributeName {
            return field.Type, nil
        }
    }
    return nil, fmt.Errorf("attribute %s not found in entity type %s", attributeName, entityName)
}

// ImplementsEntity vérifie si un type implémente l'interface Entity
func ImplementsEntity(structType reflect.Type) bool {
    entityType := reflect.TypeOf((*Entity)(nil)).Elem()
    return reflect.PtrTo(structType).Implements(entityType)
}

// FindPotentialKeyField trouve un champ potentiel qui pourrait être utilisé comme clé
func FindPotentialKeyField(entityName string) (string, reflect.Type, error) {
    fields, err := GetFieldsByEntityType(entityName)
    if err != nil {
        return "", nil, err
    }

    for _, field := range fields {
        if strings.Contains(field.Name, "Key") || strings.Contains(field.Name, "Name") {
            return field.Name, field.Type, nil
        }
    }
    return "", nil, fmt.Errorf("no potential key field found in entity type %s", entityName)
}

// NewEntityInstance crée une nouvelle instance d'une entité par son nom de type
func NewEntityInstance(typeName string) (Entity, error) {
    structType, exists := GetEntityType(typeName)
    if !exists {
        return nil, fmt.Errorf("type %s not found", typeName)
    }
    return reflect.New(structType).Interface().(Entity), nil
}