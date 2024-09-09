package models

import (
	entity "api-interface/entities"
	repository "api-interface/repositories"
)

type BucketModel struct {
    bucketRepository *repository.QueryBuilder[*entity.Bucket]
}

// NewBucketModel initialise un BucketModel avec le repository approprié
func UseBucketModel() (*BucketModel, error) {
    queryBuilder, err := UseRepository[*entity.Bucket]("Bucket")
    if err != nil {
        return nil, err
    }

    return &BucketModel{
        bucketRepository: queryBuilder,
    }, nil
}

func (bm *BucketModel) Insert(bucket *entity.Bucket) error {
	return bm.bucketRepository.Insert(bucket)
}

func (bm *BucketModel) GetAllBuckets() ([]entity.Bucket, error) {
    buckets, err := bm.bucketRepository.Find(func(b *entity.Bucket) bool {
        // Appliquer un filtre si nécessaire, par exemple, retourner true pour inclure tous les buckets
        return true
    })
    if err != nil {
        return nil, err
    }

    // Convertir []*entity.Bucket en []entity.Bucket
    var result []entity.Bucket
    for _, bucket := range buckets {
        result = append(result, *bucket)
    }
    return result, nil
}

