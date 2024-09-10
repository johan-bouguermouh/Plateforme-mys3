package controllers

import (
	entity "api-interface/entities"
	"api-interface/models"
	"fmt"

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

// Field names should start with an uppercase letter
type Person struct {
    Name string `json:"name" xml:"name" form:"name"`
    Pass string `json:"pass" xml:"pass" form:"pass"`
}
/** InsertBucket permet d'insérer un Bucket dans la base de données */
func (b *BucketController) InsertBucket(c *fiber.Ctx) error {
    // Log the raw body
    bucket := new(entity.Bucket)

    // Extraire et valider les données du formulaire
  
    if err := c.BodyParser(bucket); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
            "errors": err.Error(),
        })
    }

    //Valider les champs obligatoires
    if bucket.Name == "" || bucket.CreationDate == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "missing required bucket fields",
        })
    }

    //Insérer le bucket dans la base de données
    err := b.bucketService.Insert(bucket)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON("helloWorld")
}

//* getAllBuckets permet de récupérer tous les Buckets de la base de données
func (b *BucketController) GetAllBuckets(c *fiber.Ctx) error {
    buckets, err := b.bucketService.GetAllBuckets()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(buckets)
}

/** GetBucketByName permet de récupérer un Bucket par son nom */
func (b *BucketController) GetBucketByName(c *fiber.Ctx) error {
    name := c.Params("name")
    bucket, err := b.bucketService.GetBucketByName(name)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(bucket)
}