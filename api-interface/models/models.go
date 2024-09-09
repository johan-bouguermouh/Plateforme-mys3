package models

import (
	"fmt"
	"reflect"

	entity "api-interface/entities"
	repository "api-interface/repositories"

	"go.etcd.io/bbolt"
)

func init() {
    entity.RegisterEntityType("Bucket", reflect.TypeOf(entity.Bucket{}))
    entity.RegisterEntityType("BucketObject", reflect.TypeOf(entity.BucketObject{}))
    entity.RegisterEntityType("Owner", reflect.TypeOf(entity.Owner{}))
    registerQueryBuilder[*entity.Bucket]("Bucket")
    registerQueryBuilder[*entity.BucketObject]("BucketObject")
    registerQueryBuilder[*entity.Owner]("Owner")
}

var repositories = make(map[string]interface{})
var queryBuilderConstructors = make(map[string]func(entity.Entity, *bbolt.DB, string) (interface{}, error))

func registerQueryBuilder[T entity.Entity](typeName string) {
    queryBuilderConstructors[typeName] = func(e entity.Entity, db *bbolt.DB, bucketName string) (interface{}, error) {
        return repository.NewQueryBuilder[T](e.(T), db, bucketName)
    }
}

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

func GetRepository[T entity.Entity](entityName string) (*repository.QueryBuilder[T], error) {
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