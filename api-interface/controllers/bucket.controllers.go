package controllers

import (
	entity "api-interface/entities"
	"api-interface/models"

	"github.com/gofiber/fiber/v2"
)

/** Bucket Controller permet de gérer les requêtes liées aux Buckets */
type BucketController struct {
    bucketService *models.BucketModel
}

/** NewBucketController initialise un BucketController avec le service approprié */
func NewBucketController() (*BucketController, error) {
    bucketService, err := models.UseBucketModel()
    if err != nil {
        return nil, err
    }

    return &BucketController{
        bucketService: bucketService,
    }, nil
}

/** InsertBucket permet d'insérer un Bucket dans la base de données */
func (b *BucketController) InsertBucket(c *fiber.Ctx) error {
    println("InsertBucket", c)

    // Extraire et valider les données du formulaire
    var bucket entity.Bucket
    if err := c.BodyParser(&bucket); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "cannot parse JSON",
        })
    }

    // Valider les champs obligatoires
    if bucket.Name == "" || bucket.CreationDate == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "missing required bucket fields",
        })
    }

    // Insérer le bucket dans la base de données
    err := b.bucketService.Insert(&bucket)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(bucket)
}