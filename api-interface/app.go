package main

import (
	"flag"
	"log"

	"api-interface/models"

	"github.com/joho/godotenv"
	"go.etcd.io/bbolt"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func init() {
	// Charger les variables d'environnement
	//envPath := filepath.Join("..", ".env")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// func main() {
// 	// Parse command-line flags
// 	flag.Parse()

// 	// Connected with database
// 	database.Connect()

// 	// Create fiber app
// 	app := fiber.New(fiber.Config{
// 		Prefork: *prod, // go run app.go -prod
// 	})

// 	// Middleware
// 	app.Use(recover.New())
// 	app.Use(logger.New())

// 	// Create a /api/v1 endpoint
// 	v1 := app.Group("/api/v1")

// 	routes.Router(v1)

// 	// Bind handlers
// 	v1.Get("/users", handlers.UserList)
// 	v1.Post("/users", handlers.UserCreate)

// 	// Setup static files
// 	app.Static("/", "./static/public")

// 	// Handle not founds
// 	app.Use(handlers.NotFound)

// 	// Listen on port 3000
// 	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
// }

func main() {
	    // Ouvrir la base de données BoltDB
		db, err := bbolt.Open("my.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
	
		// Initialiser les repositories
		models.InitRepositories(db)
	
		//Optionnel : Utiliser les repositories pour vérifier que tout fonctionne
		repo := models.GetRepository("Bucket")
		if repo != nil {
			log.Println("Repository for Bucket initialized successfully")
		} else {
			log.Println("Failed to initialize repository for Bucket")
		}
	}
