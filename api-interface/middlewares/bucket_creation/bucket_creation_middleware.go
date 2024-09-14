package middlewares

import (
	entities "api-interface/entities"
	"api-interface/handlers/errors"
	bucket_name_validation "api-interface/middlewares/bucket_name_validation"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// BucketValidationMiddleware est le middleware de validation pour les buckets
func BucketValidationMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Lire le corps de la requête
        bodyRequest := c.Body()
        bucketName := c.Params("bucketName")


        // Nouvelle instance de CreateBucketRequestStruct
        bucketRequest := new(entities.CreateBucketRequestStruct)
        if bodyRequest == nil {
        // Décoder le corps XML dans la struct bucketRequest
        if err := xml.Unmarshal(bodyRequest, bucketRequest); err != nil {
            return errors.HandleError(c, errors.ErrUnprocessableEntity, "Erreur de parsing XML: " + err.Error())
        }

        // Valider les champs du bucket
        if err := validateBucketFields(bucketRequest, c); err != nil {
            return err
        }

        c.Locals("bucketRequest", bucketRequest)
    }


        // Créez un validateur pour le nom du bucket
        validator := bucket_name_validation.NewBucketNameValidator()
        result := validator.Validate(bucketName)

        if len(result) > 0 {
            errorMessage := strings.Join(result, "; ")
            return errors.HandleError(c, errors.ErrBadRequest, errorMessage)
        }

        // Vérifier si le bucket existe déjà
        if bucketExists(bucketName) {
            return errors.HandleError(c, errors.ErrBadRequest, "Le bucket avec ce nom existe déjà.")
        }

        // Attacher les données validées au contexte
        c.Locals("bucketName", bucketName)

        return c.Next()
    }
}

// validateBucketFields valide les champs du bucket
func validateBucketFields(bucket *entities.CreateBucketRequestStruct, c *fiber.Ctx) error {
    if bucket.Owner == (entities.Owner{}) {
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

// bucketExists vérifie si un dossier de bucket existe déjà
func bucketExists(bucketName string) bool {
    bucketPath := filepath.Join("./buckets", bucketName)
    info, err := os.Stat(bucketPath)
    return !os.IsNotExist(err) && info.IsDir()
}
