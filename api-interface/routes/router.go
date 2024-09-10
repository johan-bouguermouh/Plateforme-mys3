package routes

import (
	"api-interface/database"
	"api-interface/handlers"

	"github.com/gofiber/fiber/v2"
)

// Router configure les routes de l'application
func Router(app *fiber.App) {
	// Routes publiques
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Bienvenue sur la plateforme MyS3")
	// })
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("static/public/index.html")
	})

	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Routes protégées
	protected := app.Group("")
	protected.Use(handlers.AuthRequired)
	protected.Post("/bucket", database.CreateBucket)
	protected.Get("/bucket/:bucketName/files", database.ListFiles)
	protected.Post("/bucket/:bucketName/upload", database.UploadFile)
	protected.Get("/bucket/:bucketName/file/:fileName", database.DownloadFile)
	protected.Delete("/bucket/:bucketName/file/:fileName", database.DeleteFile)
}
