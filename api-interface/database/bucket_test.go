package database

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateBucket(t *testing.T) {
	app := fiber.New()
	app.Post("/bucket", CreateBucket)

	req := httptest.NewRequest(http.MethodPost, "/bucket", bytes.NewBufferString(`{"bucketName":"testBucket"}`))
	req.Header.Set("Content-Type", "application/json")

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
	assert.Equal(t, "Bucket créé avec succès", responseBody)
}
