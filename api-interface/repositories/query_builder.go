package repository

import (
	entity "api-interface/entities"
	"fmt"

	"go.etcd.io/bbolt"
)

// QueryBuilder est une structure générique qui implémente les opérations CRUD pour une entité donnée. Cette Factory exploite les entité implémentant l'interface Entity pour effectuer des opérations de base de données.
type QueryBuilder[T entity.Entity] struct {
    entity     T
    db         *bbolt.DB
    bucketName string
}

/** Retourne l'entité associée au QueryBuilder
    * @return T, l'entité associée au QueryBuilder
*/
func (qb *QueryBuilder[T]) GetEntity() T {
    return qb.entity
}

/** Crée une nouvelle instance de QueryBuilder pour une entité donnée
    * @return string, le nom du bucket associé au QueryBuilder
*/
func NewQueryBuilder[T entity.Entity](entity T, db *bbolt.DB, bucketName string) (*QueryBuilder[T], error)  {
    // Vérifier et créer le bucket si nécessaire
    err := db.Update(func(tx *bbolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte(bucketName))
        if err != nil {
            return fmt.Errorf("erreur lors de la création du bucket %s : %v", bucketName, err)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return &QueryBuilder[T]{entity: entity, db: db, bucketName: bucketName}, nil
}

/** Provate : Retourne le bucket associé au QueryBuilder
    * @return *bbolt.Bucket, le bucket associé au QueryBuilder
*/
func (qb *QueryBuilder[T]) getBucket(tx *bbolt.Tx) (*bbolt.Bucket, error) {
    bucket := tx.Bucket([]byte(qb.bucketName))
    if bucket == nil {
        return nil, fmt.Errorf("bucket %s not found", qb.bucketName)
    }
    return bucket, nil
}

/** Insère une entité dans la base de données
    * @param entity T, l'entité à insérer
    * @return error, une erreur si l'opération a échoué
*/
func (qb *QueryBuilder[T]) Insert(entity T) error {
    return qb.db.Update(func(tx *bbolt.Tx) error {
        bucket, err := qb.getBucket(tx)
        if err != nil {
            return err
        }
        data, err := entity.Serialize()
        if err != nil {
            return err
        }
        return bucket.Put([]byte(entity.GetKey()), data)
    })
}

/** Met à jour une entité dans la base de données
    * @param entity T, l'entité à mettre à jour
    * @return error, une erreur si l'opération a échoué
*/
func (qb *QueryBuilder[T]) Get(key string, entity T) error {
    return qb.db.View(func(tx *bbolt.Tx) error {
        bucket, errBucket := qb.getBucket(tx)
        if errBucket != nil {
            return errBucket
        }
        data := bucket.Get([]byte(key))
        if data == nil {
            return fmt.Errorf("entity not found")
        }
        err := entity.Deserialize(data)
        if err != nil {
            return err
        }
        return nil
    })
}

/** Met à jour une entité dans la base de données
    * @param entity T, l'entité à mettre à jour
    * @return error, une erreur si l'opération a échoué
*/
func (qb *QueryBuilder[T]) GetAll(filter func(T) bool) ([]T, error) {
    var entities []T
    err := qb.db.View(func(tx *bbolt.Tx) error {
        bucket := tx.Bucket([]byte(qb.bucketName))
        return bucket.ForEach(func(k, v []byte) error {
            var entity T
            if err := entity.Deserialize(v); err != nil {
                return err
            }
            if filter(entity) {
                entities = append(entities, entity)
            }
            return nil
        })
    })
    return entities, err
}