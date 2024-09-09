package models

import (
	"fmt"
	"reflect"

	entity "api-interface/entities"
	repository "api-interface/repositories"

	"go.etcd.io/bbolt"
)

func init() {
    Use[*entity.Bucket]()
    Use[*entity.BucketObject]()
    Use[*entity.Owner]()
}

var repositories = make(map[string]interface{})
var queryBuilderConstructors = make(map[string]func(entity.Entity, *bbolt.DB, string) (interface{}, error))

/** Déclare un type d'entité et initialise son QueryBuilder, le Type doit être un pointeur vers une structure qui implémente l'interface entity.Entity, et doit être déclaré dans le package entities, la méthode prend en paramètre le type de la structure définie en tant qu'object reflect.Type
    * @Type T entity.Entity 
    * @ATTENTION : Le type d'entité doit être un pointeur vers une structure qui implémente l'interface Entity
*/
func Use[T entity.Entity]() {
    typeName, instanceType := getStructType[T]()
    entity.RegisterEntityType(typeName, instanceType)
    queryBuilderConstructors[typeName] = func(e entity.Entity, db *bbolt.DB, bucketName string) (interface{}, error) {
        return repository.NewQueryBuilder[T](e.(T), db, bucketName)
    }
}

// Private | Retourne une instance de l'entité et son type selon son pointer
func getStructType[T entity.Entity]() (string, reflect.Type) {
    var instance T
    instanceType := reflect.TypeOf(instance).Elem()
    typeName := instanceType.Name()
    return typeName, instanceType
}

/** Initialise les repositories pour chaque entité déclarée utilisé dans la fonction init() de ce fichier. Cette fonction est appelée au démarrage de l'application pour initialiser les repositories pour chaque entité déclarée. Init repositories utilise les QueryBuilder pour chaque entité déclarée et les stocke dans un map pour une utilisation ultérieure.
    * @param db *bbolt.DB, instance de la base de données bolt
    @ATTENTION : Cette fonction doit être appelée au démarrage de l'application pour initialiser les repositories pour chaque entité déclarée. Il est conseillé de l'appeler dans la fonction main() de l'application pour garantir que les repositories sont initialisés avant toute utilisation.
*/
func InitRepositories(db *bbolt.DB) {
    fmt.Println(Yellow + "=== Initialisation des Repositories ===" + Reset)

    for typeName, structType := range entity.GetAllEntityTypes() {
        if AssureEntityCompliance(typeName, structType) {
            bucketName := structType.Name() + "s"
            entityInstance, err := entity.NewEntityInstance(typeName)
            if err != nil {
                fmt.Printf("%sErreur lors de la création de l'instance de l'entité %s : %s%s\n", Red, typeName, err, Reset)
                continue
            }

            constructor, exists := queryBuilderConstructors[typeName]
            if !exists {
                fmt.Printf("%sType d'entité non supporté : %s%s\n", Red, typeName, Reset)
                continue
            }

            queryBuilder, err := constructor(entityInstance, db, bucketName)
            if err != nil {
                fmt.Printf("%sErreur lors de la création du QueryBuilder pour l'entité %s : %s%s\n", Red, typeName, err, Reset)
                continue
            }

            repositories[structType.Name()] = queryBuilder
        } else {
            fmt.Printf("%sL'entité %s n'implémente pas l'interface entity.Entity%s\n", Red, structType.Name(), Reset)
        }
    }

    fmt.Println(Yellow + "=== Fin de l'initialisation des Repositories ===" + Reset)
}

/** Instancie un Repository pour une entité donnée. Cette fonction peut être appeler dans le init() ou dans la fonction d'instanctiation du modèle d'un modèle. Elle récupère l'instance du QueryBuilder pour une entité donnée, ce dernier est directeent liées a un bucket spécifique lié à l'entité.
    * @Type T entity.Entity
    * @param entityName string
    * @return *repository.QueryBuilder[T]
    * @return error
*/
func UseRepository[T entity.Entity](entityName string) (*repository.QueryBuilder[T], error) {
    repo, ok := repositories[entityName]
    if !ok {
        return nil, fmt.Errorf("repository for %s not found", entityName)
    }

    queryBuilder, ok := repo.(*repository.QueryBuilder[T])
    if !ok {
        return nil, fmt.Errorf("repository is not of type *repositories.QueryBuilder[%s]", entityName)
    }

    return queryBuilder, nil
}