package models

import (
	"fmt"

	entity "api-interface/entities"
	repository "api-interface/repositories"

	"go.etcd.io/bbolt"
)

var (
    repositories = make(map[string]interface{})
)

func InitRepositories(db *bbolt.DB) {
    fmt.Println(Yellow + "=== Initialisation des Repositories ===" + Reset)

    // Obtenir tous les types enregistrés dans le registre global
    for typeName, structType := range entity.GetAllEntityTypes() {

        if AssureEntityCompliance(typeName, structType) {
            bucketName := structType.Name() + "s"
            // fmt.Printf("Nom du bucket : %s\n", bucketName)
            entityInstance, err := entity.NewEntityInstance(typeName)
            if err != nil {
                fmt.Printf("%sErreur lors de la création de l'instance de l'entité %s : %s%s\n", Red, typeName, err, Reset)
                continue
            }

            // Utiliser un switch pour gérer les types d'entités spécifiques
            var queryBuilder interface{}
            switch e := entityInstance.(type) {
            case *entity.Bucket:
                // on retourn e est on l'exprime en tant que *entity.Bucket
                fmt.Println("Traitement de l'entité reconnue comme Bucket :", e)
                queryBuilder, err = repository.NewQueryBuilder[*entity.Bucket](entityInstance.(*entity.Bucket), db, bucketName)

            case *entity.BucketObject:
                fmt.Println("Traitement de l'entité comme BucketObject :", e)
                queryBuilder, err = repository.NewQueryBuilder[*entity.BucketObject](entityInstance.(*entity.BucketObject), db, bucketName)
            case *entity.Owner:
                fmt.Println("Traitement de l'entité comme Owner:", e)
                queryBuilder, err = repository.NewQueryBuilder[*entity.Owner](entityInstance.(*entity.Owner), db, bucketName)
            default:
                fmt.Printf("%sType d'entité non supporté : %s%s\n", Red, typeName, Reset)
                continue
            }
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