package repository

import (
	entity "api-interface/entities"
	"fmt"

	"go.etcd.io/bbolt"
)

type QueryBuilder[T entity.Entity] struct {
    entity     T
    db         *bbolt.DB
    bucketName string
}

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

// Méthode privée pour accéder au bucket
func (qb *QueryBuilder[T]) getBucket(tx *bbolt.Tx) (*bbolt.Bucket, error) {
    bucket := tx.Bucket([]byte(qb.bucketName))
    if bucket == nil {
        return nil, fmt.Errorf("bucket %s not found", qb.bucketName)
    }
    return bucket, nil
}

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