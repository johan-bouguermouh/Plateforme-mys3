package database

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Liste des fichiers dans un bucket
func ListFiles(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	files, err := os.ReadDir("./data/" + bucketName)
	if err != nil {
		return c.Status(500).SendString("Erreur lors de la lecture du bucket")
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return c.JSON(fileNames)
}

// Upload d'un fichier

func UploadFile(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	bucketPath := "./data/" + bucketName

	// Vérifier si le bucket existe
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return c.Status(400).SendString("Le bucket n'existe pas")
	}

	// Récupérer le fichier depuis la requête
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Erreur lors de la réception du fichier : %v", err)
		return c.Status(400).SendString("Fichier requis")
	}

	// Créer un chemin complet pour le fichier
	filePath := bucketPath + "/" + file.Filename

	// Sauvegarder le fichier sur le disque
	err = c.SaveFile(file, filePath)
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde du fichier : %v", err)
		return c.Status(500).SendString("Erreur lors de l'upload du fichier")
	}

	return c.SendString("Fichier uploadé avec succès")
}

// Téléchargement d'un fichier
func DownloadFile(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	fileName := c.Params("fileName")

	filePath := "./data/" + bucketName + "/" + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).SendString("Fichier non trouvé")
	}

	return c.SendFile(filePath)
}

// Suppression d'un fichier
func DeleteFile(c *fiber.Ctx) error {
	bucketName := c.Params("bucketName")
	fileName := c.Params("fileName")

	filePath := "./data/" + bucketName + "/" + fileName
	log.Printf("Suppression du fichier : %s", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("Fichier non trouvé")
		return c.Status(404).SendString("Fichier non trouvé")
	}

	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Erreur lors de la suppression du fichier : %v", err)
		return c.Status(500).SendString("Erreur lors de la suppression du fichier")
	}

	return c.SendString("Fichier supprimé avec succès")
}
