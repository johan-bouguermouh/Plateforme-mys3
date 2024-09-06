package routes

import (
	"api-interface/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(route fiber.Router) {
    BucketController, err := controllers.NewBucketController()
    if err != nil {
        panic(err)
    }
    /** Route for Minio Upload */
    //route.Post("/upload", controllers.UploadFile)
    //création d'un sous-groupe pour les routes de l'API Bucket
    route.Group("/bucket")
    {
        route.Post("/", BucketController.InsertBucket)
    }
}