package database

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	app := fiber.New()
	app.Get("/bucket/:bucketName/files", ListFiles)

	// Créer un bucket de test
	bucketName := "testBucket"
	os.Mkdir("./data/"+bucketName, 0755)
	defer os.RemoveAll("./data/" + bucketName) // Nettoyage

	fileName := "testFile.txt"
	err := os.WriteFile("./data/"+bucketName+"/"+fileName, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier de test : %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/bucket/"+bucketName+"/files", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Erreur lors du test : %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erreur lors de la lecture du corps de la réponse : %v", err)
	}
	defer resp.Body.Close()

	var fileNames []string
	err = json.Unmarshal(body, &fileNames)
	if err != nil {
		t.Fatalf("Erreur lors de la décompression de la réponse : %v", err)
	}

	assert.Contains(t, fileNames, fileName)
}

func TestUploadFile(t *testing.T) {
	app := fiber.New()
	app.Post("/bucket/:bucketName/upload", UploadFile)

	// Créer un bucket de test
	bucketName := "testBucket"
	os.Mkdir("./data/"+bucketName, 0755)
	defer os.RemoveAll("./data/" + bucketName) // Nettoyage

	// Créer un fichier multipart pour l'upload
	file := &bytes.Buffer{}
	writer := multipart.NewWriter(file)
	part, err := writer.CreateFormFile("file", "testFile.txt")
	if err != nil {
		t.Fatalf("Erreur lors de la création du fichier multipart : %v", err)
	}
	part.Write([]byte("test content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/bucket/"+bucketName+"/upload", file)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Erreur lors du test : %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Erreur lors de la lecture du corps de la réponse : %v", err)
	}
	defer resp.Body.Close()

	responseBody := string(body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Fichier uploadé avec succès", responseBody)
}
