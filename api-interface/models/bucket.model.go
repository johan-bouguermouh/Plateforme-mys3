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
    queryBuilder, err := GetRepository[*entity.Bucket]("Bucket")
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

