package routes

import (
	"api-interface/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(route fiber.Router) {
    /** Route for Minio Upload */
    route.Post("/upload", controllers.UploadFile)
}