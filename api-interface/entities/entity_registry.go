package entity

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

/** EntityRegistry est une structure qui contient une carte des types d'entités enregistrés.
 * @mu Défini mutex pour la synchronisation
 * @types Défini une carte des types d'entités
 */
type EntityRegistry struct {
    mu    sync.RWMutex
    types map[string]reflect.Type
}

/** FieldInfo contient des informations sur un champ d'une entité
    * @Name Défini le nom du champ, simplifie l'accès aux champs de la structure
    * @Type Défini le type du champ, type de la structure de réflexion (Object)
    * @Tag Défini les tags du champ, enregistrement des métadonnées structurées
*/
type FieldInfo struct {
    Name string
    Type reflect.Type
    Tag  reflect.StructTag
}

/** Crée une nouvelle instance de EntityRegistry. L'instance EntityRegistry est utilisée pour enregistrer et récupérer des types d'entités. Sa création permet de gérer les types d'entités de manière centralisée.
    * @return *EntityRegistry, un pointeur vers une instance de EntityRegistry
*/
func NewEntityRegistry() *EntityRegistry {
    return &EntityRegistry{
        types: make(map[string]reflect.Type),
    }
}

/** Enregistre un type d'entité dans une instance spécifique de EntityRegistry. En utilisant plusieurs instances, vous pouvez répartir la charge entre différentes instances, ce qui peut améliorer la performance et la scalabilité. Cependant, cela nécessite une gestion plus complexe des instances et de la synchronisation.
    * @portée : Méthode d'instance de la structure EntityRegistry.
    * @Accès : Appelée sur une instance de EntityRegistry.
    * @param name string, le nom du type d'entité
    * @param t reflect.Type, le type d'entité à enregistrer
*/
func (r *EntityRegistry) Register(name string, t reflect.Type) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.types[name] = t
}

/** Récupère un type d'entité à partir de la carte. Cette méthode est utilisée pour récupérer un type d'entité enregistré dans le registre. Si le type d'entité n'existe pas, la méthode renvoie false.
    * @portée : Méthode d'instance de la structure EntityRegistry.
    * @Accès : Appelée sur une instance de EntityRegistry.
    * @param name string, le nom du type d'entité
    * @return reflect.Type, bool, le type d'entité et un booléen indiquant si le type d'entité a été trouvé
*/
func (r *EntityRegistry) Get(name string) (reflect.Type, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    t, ok := r.types[name]
    return t, ok
}

// Créer une nouvelle instance de EntityRegistry
var globalRegistry = NewEntityRegistry()

/** Enregistrer un type d'entité dans l'instance globale globalRegistry de EntityRegistry. Pour une application massive, il est souvent préférable de commencer avec une approche centralisée et de passer à une approche distribuée au fur et à mesure que les besoins en scalabilité augmentent. Attention cependant, instance globale peut devenir un goulot d'étranglement si de nombreuses goroutines tentent d'y accéder simultanément.
    * Portée : Fonction globale.
    * Accès : Appelée directement sans avoir besoin d'une instance de EntityRegistry.
    * @param name string, le nom du type d'entité
    * @param t reflect.Type, le type d'entité à enregistrer
    * @return void
*/
func RegisterEntityType(name string, t reflect.Type) {
    globalRegistry.Register(name, t)
}

/** GetEntityType récupère un type d'entité à partir de la carte globale. 
    * Portée : Fonction propre au registre global.
    * Accès : Appelée directement sans avoir besoin d'une instance de EntityRegistry.
    * @param name string, le nom du type d'entité
    * @return reflect.Type, bool, le type d'entité et un booléen indiquant si le type d'entité a été trouvé
*/
func GetEntityType(name string) (reflect.Type, bool) {
    return globalRegistry.Get(name)
}

/** Retourne tous les types d'entités enregistrés dans le registre global.
    * Portée : Fonction propre au registre global.
    * Accès : Appelée directement sans avoir besoin d'une instance de EntityRegistry.
    * @return map[string]reflect.Type, une carte des types d'entités
*/
func GetAllEntityTypes() map[string]reflect.Type {
    globalRegistry.mu.RLock()
    defer globalRegistry.mu.RUnlock()
    return globalRegistry.types
}

/** GetFieldsByEntityType retourne tous les attributs présents dans une entité
    * @param name string, le nom du type d'entité
    * @return []FieldInfo, error, une liste des attributs structurés par Name, Type et Tag, présent dans le registre global, retourne erreur si le type d'entité n'est pas trouvé
*/
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

/** Permet de vérifier si un attribut existe dans une entité
    * @param entityName string, le nom du type d'entité
    * @param attributeName string, le nom de l'attribut à vérifier
    * @return bool, error, un booléen indiquant si l'attribut existe et une erreur si le type d'entité n'est pas trouvé
*/
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

/** Retourne le type d'un attribut d'une entité par son nom
    * @param entityName string, le nom du type d'entité
    * @param attributeName string, le nom de l'attribut
    * @return reflect.Type, error, le type de l'attribut et une erreur si le type d'entité ou l'attribut n'est pas trouvé
*/
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

/** Vérifie si un type d'entité implémente l'interface Entity. Est utile lors de l'intégration d'un Repository pour s'assurer que le type d'entité est conforme.
    * @param structType reflect.Type, le type d'entité à vérifier
    * @return bool, un booléen indiquant si le type d'entité implémente l'interface Entity
*/
func ImplementsEntity(structType reflect.Type) bool {
    entityType := reflect.TypeOf((*Entity)(nil)).Elem()
    return reflect.PtrTo(structType).Implements(entityType)
}

/** Fonction propre à la création d'un Repository, permet de trouver un champ potentiel de clé dans une entité. Un champ potentiel de clé est un champ qui contient le mot "Key" ou "Name" dans son nom. Cette fonction est une aide au développeur pour respecter des règles de nommage cohérentes lors la création d'un bucket bolt.
    * @param entityName string, le nom du type d'entité
    * @return string, reflect.Type, error, le nom du champ potentiel de clé, le type du champ potentiel de clé et une erreur si aucun champ potentiel de clé n'est trouvé.
*/
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

/** Crée une nouvelle instance d'une entité par son nom de type, cette approche est utilisée pour instancier une entité sans connaître son type à l'avance. Permet ainsi, de créer une instance lors de l'usage du design pattern Factory (pour exemple).
    * @param typeName string, le nom du type d'entité
    * @return Entity, error, une instance de l'entité et une erreur si le type d'entité n'est pas trouvé
    * @Attention : Cette fonction utilise la réflexion pour créer une instance d'une entité, ce qui peut entraîner une perte de performance. Il est recommandé d'utiliser cette fonction uniquement lorsqu'il est impossible de connaître le type d'entité à l'avance.
*/
func NewEntityInstance(typeName string) (Entity, error) {
    structType, exists := GetEntityType(typeName)
    if !exists {
        return nil, fmt.Errorf("type %s not found", typeName)
    }
    return reflect.New(structType).Interface().(Entity), nil
}