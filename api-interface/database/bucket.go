package database

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Création d'un bucket
// func CreateBucket(c *fiber.Ctx) error {
// 	type BucketRequest struct {
// 		BucketName string `json:"bucketName"`
// 	}

// 	req := new(BucketRequest)
// 	if err := c.BodyParser(req); err != nil {
// 		log.Printf("Erreur de parsing : %v", err)
// 		return c.Status(400).SendString("Données invalides")
// 	}

// 	bucketName := req.BucketName
// 	if bucketName == "" {
// 		log.Println("Le nom du bucket est vide")
// 		return c.Status(400).SendString("Le nom du bucket est requis")
// 	}

// 	bucketPath := "/data/" + bucketName
// 	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
// 		return c.Status(400).SendString("Le bucket existe déjà")
// 	}

// 	dataPath := "/data/"
// 	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
// 		err := os.Mkdir(dataPath, 0755)
// 		if err != nil {
// 			log.Printf("Erreur lors de la création du répertoire de données : %v", err)
// 			return c.Status(500).SendString("Erreur lors de la configuration du serveur")
// 		}
// 	}

//		return c.SendString("Bucket créé avec succès")
//	}
//
// Création d'un bucket
func CreateBucket(c *fiber.Ctx) error {
	type BucketRequest struct {
		BucketName string `json:"bucketName"`
	}

	req := new(BucketRequest)
	if err := c.BodyParser(req); err != nil {
		log.Printf("Erreur de parsing : %v", err)
		return c.Status(400).SendString("Données invalides")
	}

	bucketName := req.BucketName
	if bucketName == "" {
		log.Println("Le nom du bucket est vide")
		return c.Status(400).SendString("Le nom du bucket est requis")
	}

	// Chemin du bucket
	bucketPath := "./data/" + bucketName

	// Vérifier si le bucket existe déjà
	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
		return c.Status(400).SendString("Le bucket existe déjà")
	}

	// Créer le dossier /data/ si nécessaire
	if _, err := os.Stat("./data/"); os.IsNotExist(err) {
		err := os.Mkdir("./data/", 0755)
		if err != nil {
			log.Printf("Erreur lors de la création du répertoire de données : %v", err)
			return c.Status(500).SendString("Erreur lors de la configuration du serveur")
		}
	}

	// Créer le bucket
	err := os.Mkdir(bucketPath, 0755)
	if err != nil {
		log.Printf("Erreur lors de la création du bucket : %v", err)
		return c.Status(500).SendString("Erreur lors de la création du bucket")
	}

	return c.SendString("Bucket créé avec succès")
}
