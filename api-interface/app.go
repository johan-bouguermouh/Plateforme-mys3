// main.go
package main

import (
	"api-interface/database"
	"api-interface/handlers"
	"api-interface/models"
	"api-interface/routes"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"go.etcd.io/bbolt"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func init() {
	// Charger les variables d'environnement
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	// Parse command-line flags
	flag.Parse()

	// Connected with database
	db, err := bbolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repositories
	models.InitRepositories(db)
	
	// Connect with database
	database.Connect()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Routes
	routes.Router(app)

	// Bind handlers
	app.Get("/users", handlers.UserList)
	app.Post("/users", handlers.UserCreate)

	// Setup static files
	app.Static("/", "./static/public")

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port
	log.Fatal(app.Listen(":3000"))
}