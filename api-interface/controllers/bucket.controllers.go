package controllers

import (
    utils "api-interface/utils"
    entities "api-interface/entities"
    "api-interface/models"
    "api-interface/handlers/errors"

    "encoding/xml"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "time"
    "os"
    "path/filepath"
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

func (b *BucketController) InsertBucket(c *fiber.Ctx) error {
    // Lire le corps de la requête
    bodyRequest := c.Body()

    // Nouvelle instance de CreateBucketRequestStruct (struct utilisée pour recevoir les données de la requête)
    bucketRequest := new(entities.CreateBucketRequestStruct)

    // Décoder le corps XML dans la struct bucketRequest
    if err := xml.Unmarshal(bodyRequest, bucketRequest); err != nil {
        return errors.HandleError(c, errors.ErrUnprocessableEntity, "Erreur de parsing XML: " + err.Error())
    }

    // Valider les champs du bucket (à partir de la struct de la requête)
    if err := validateBucketFields(bucketRequest, c); err != nil {
        // Si la validation échoue, la fonction `validateBucketFields` renverra une réponse appropriée
        return err
    }

    // Vérifier si le bucket existe déjà
    if bucketExists(bucketRequest.Name) {
        return errors.HandleError(c, errors.ErrBadRequest, "Le bucket avec ce nom existe déjà.")
    }

    // Générer l'URI pour le bucket (basé sur son nom ou d'autres critères)
    bucketURI := generateBucketURI(bucketRequest.Name)

    // Créer un répertoire pour le bucket (émulation du S3 bucket)
    bucketPath, err := createBucketDirectory(bucketRequest.Name)
    if err != nil {
        return errors.HandleError(c, errors.ErrInternalServerError, "Erreur lors de la création du répertoire du bucket")
    }

    defaultBucketCreationDate := time.Now()

    // Conversion de CreateBucketRequestStruct vers Bucket (l'entité qui sera insérée en base de données)
    bucket := &entities.Bucket{
        Name:         bucketRequest.Name,
        CreationDate: utils.StringPointer(defaultBucketCreationDate.String()),
        Owner:        bucketRequest.Owner,
        URI:          bucketURI, // URI généré par le serveur
        Type:         bucketRequest.Type,
        Versioning:   bucketRequest.Versioning,
    }

    // Insérer le bucket dans la base de données via le service
    err = b.bucketService.Insert(bucket)
    if err != nil {
        return errors.HandleError(c, errors.ErrInternalServerError, "Erreur lors de l'insertion du bucket dans la base de données")
    }

    // Réponse réussie avec les informations du bucket
    return c.JSON(fiber.Map{
        "message": "Bucket inséré avec succès",
        "uri":     bucketURI,
        "path":    bucketPath, // Le chemin du répertoire du bucket sur le serveur
    })
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

// Fonction de validation des champs du Bucket
func validateBucketFields(bucket *entities.CreateBucketRequestStruct, c *fiber.Ctx) error { // Utiliser entities.Bucket ici

    // Champs obligatoires 
    if bucket.Name == "" {
        return errors.HandleError(c, errors.ErrBadRequest, "Le nom du bucket est obligatoire.")
    }

    if bucket.Owner == (entities.Owner{}) { // Utiliser entities.Owner ici
        return errors.HandleError(c, errors.ErrBadRequest, "Le bucket doit posséder un propriétaire.")
    }

    if bucket.Type == "" {
        return errors.HandleError(c, errors.ErrBadRequest, "Le type de bucket doit être précisé.")
    }
    
    if bucket.Versioning == "" {
        return errors.HandleError(c, errors.ErrBadRequest, "Pas d'état de version pour le bucket.")
    }

    return nil
}

func generateBucketURI(bucketName string) string {
    baseURL := "http://localhost:3000/buckets" // L'URL de base de ton application
    return fmt.Sprintf("%s/%s", baseURL, bucketName)
}

// Fonction pour vérifier si un dossier existe déjà
func bucketExists(bucketName string) bool {
    bucketPath := filepath.Join("./buckets", bucketName) // Spécifiez le chemin où les buckets sont stockés
    info, err := os.Stat(bucketPath)
    if os.IsNotExist(err) {
        return false
    }
    return info.IsDir()
}

// Créer un répertoire pour le bucket
func createBucketDirectory(bucketName string) (string, error) {
    // Définir l'emplacement du répertoire de stockage
    basePath := "./buckets" // Cela pourrait être un chemin absolu sur ton serveur

    // Construire le chemin complet pour le nouveau bucket
    bucketPath := filepath.Join(basePath, bucketName)

    // Créer le répertoire pour le bucket
    if err := os.MkdirAll(bucketPath, os.ModePerm); err != nil {
        return "", fmt.Errorf("Erreur lors de la création du bucket : %v", err)
    }

    return bucketPath, nil
}