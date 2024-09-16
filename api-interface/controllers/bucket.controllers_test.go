package controllers

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	config "api-interface/entities/bucketDTOs"

	"github.com/gofiber/fiber/v2"
)

// Création de Test pour le controller Bucket
func TestBucketController(t *testing.T) {
	// On verfie que l'instatiation de notre bucket soit correcte
	t.Run("TestNewBucketController", func(t *testing.T) {
		// On instancie le controller
		bucketController, err := NewBucketController()
		if err != nil {
			t.Fatalf("Erreur lors de l'instanciation du controller : %v", err)
		}
		if bucketController == nil {
			t.Fatalf("Erreur lors de l'instanciation du controller : %v", err)
		}
		//On verfie que le service est bien instancié
		if bucketController.bucketService == nil {
			t.Fatalf("Erreur lors de l'instanciation du service")
		}
	})

	// On verifie que l'insertion d'un bucket se passe bien
	t.Run("TestInsertBucket", func(t *testing.T) {
		app := fiber.New()
		// On instancie le controller
		bucketController, err := NewBucketController()
		if err != nil {
			t.Fatalf("Erreur lors de l'instanciation du controller : %v", err)
		}
		if bucketController == nil {
			t.Fatalf("Erreur lors de l'instanciation du controller : %v", err)
		}
		locationConstraint := "EU"
		// On construit un BOdy en XML qui correspond à la demande
		CreateBucketConfiguration := config.CreateBucketConfiguration{
			XMLNS: "http://s3.amazonaws.com/doc/2006-03-01/",
			LocationConstraint: &locationConstraint,
			Location: &config.LocationInfo{
				Name: "EU",
				Type: "AvailabilityZone",
			},
			Bucket: &config.BucketInfo{
				DataRedundancy: "SingleAvailabilityZone",
				Type:           "Directory",
			},
		}
		// on %Marshal le body en XML
		body, err := xml.Marshal(CreateBucketConfiguration)
		if err != nil {
			t.Fatalf("Erreur lors de la création du body : %v", err)
		}

		bodyReader := bytes.NewReader(body)


		// On verifie que la requête se passe bien
		req := httptest.NewRequest(http.MethodPut, "/UniqueBucket", bodyReader)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Erreur lors du test : %v", err)
		}
		// On verifie que le status de la réponse est bien 200
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Erreur lors de la création du bucket : %v", err)
		}
	})
	}
