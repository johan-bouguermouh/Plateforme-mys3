package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// Test du middleware ValidateDirectUpload

func TestValidateDirectUpload_Success(t *testing.T) {
	// Création d'un fichier de 4 Go, donc valide pour l'upload direct (moins de 5 Go).
	fileSize := 4 * 1024 * 1024 * 1024
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(strings.Repeat("a", fileSize)))
	req.Header.Set("Content-Length", strconv.Itoa(fileSize))

	// Création d'un ResponseRecorder pour capturer la réponse
	rr := httptest.NewRecorder()

	// Utilisation du middleware ValidateDirectUpload avec un handler qui renvoie 200 OK.
	handler := ValidateDirectUpload(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusOK {
		t.Errorf("Code de statut attendu %v, mais obtenu %v", http.StatusOK, rr.Code)
	}
}

func TestValidateDirectUpload_FileTooLarge(t *testing.T) {
	// Création d'un fichier de 6 Go, ce qui dépasse la limite de 5 Go pour l'upload direct.
	fileSize := 6 * 1024 * 1024 * 1024
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(strings.Repeat("a", fileSize)))
	req.Header.Set("Content-Length", strconv.Itoa(fileSize))

	// Création d'un ResponseRecorder pour capturer la réponse
	rr := httptest.NewRecorder()

	// Utilisation du middleware ValidateDirectUpload avec un handler dummy.
	handler := ValidateDirectUpload(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Vérification que l'erreur est bien renvoyée
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Code de statut attendu %v, mais obtenu %v", http.StatusRequestEntityTooLarge, rr.Code)
	}

	expectedErrorMessage := ErrFileTooLarge.Error()
	receivedErrorMessage := strings.TrimSpace(rr.Body.String())
	if receivedErrorMessage != expectedErrorMessage {
		t.Errorf("Message d'erreur attendu %v, mais obtenu %v", expectedErrorMessage, receivedErrorMessage)
	}
}

func TestValidateDirectUpload_InvalidContentLength(t *testing.T) {
	// Création d'une requête avec une longueur de contenu invalide
	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader("dummy content"))
	req.Header.Set("Content-Length", "invalid_length")

	// Création d'un ResponseRecorder pour capturer la réponse
	rr := httptest.NewRecorder()

	// Utilisation du middleware ValidateDirectUpload avec un handler dummy.
	handler := ValidateDirectUpload(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Vérification que l'erreur est bien renvoyée
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Code de statut attendu %v, mais obtenu %v", http.StatusBadRequest, rr.Code)
	}
}

// Test du middleware ValidateMultipartUpload

func TestValidateMultipartUpload_Success(t *testing.T) {
	// Simule un multipart upload avec des tailles de parties valides
	req := httptest.NewRequest(http.MethodPost, "/multipart-upload", nil)

	// Création d'un ResponseRecorder pour capturer la réponse
	rr := httptest.NewRecorder()

	// Utilisation du middleware ValidateMultipartUpload avec un handler qui renvoie 200 OK.
	handler := ValidateMultipartUpload(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusOK {
		t.Errorf("Code de statut attendu %v, mais obtenu %v", http.StatusOK, rr.Code)
	}
}

func TestValidateMultipartUpload_PartTooSmall(t *testing.T) {
	// Simule un multipart upload avec une partie plus petite que la taille minimale de 5 Mo
	extractPartSizes = func(r *http.Request) []int64 {
		return []int64{2 * 1024 * 1024, 6 * 1024 * 1024, 10 * 1024 * 1024} // La première partie est inférieure à 5 Mo
	}

	req := httptest.NewRequest(http.MethodPost, "/multipart-upload", nil)

	// Création d'un ResponseRecorder pour capturer la réponse
	rr := httptest.NewRecorder()

	// Utilisation du middleware ValidateMultipartUpload
	handler := ValidateMultipartUpload(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Vérification que l'erreur est bien renvoyée
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Code de statut attendu %v, mais obtenu %v", http.StatusRequestEntityTooLarge, rr.Code)
	}

	expectedErrorMessage := ErrPartTooSmall.Error()
	receivedErrorMessage := strings.TrimSpace(rr.Body.String())
	if receivedErrorMessage != expectedErrorMessage {
		t.Errorf("Message d'erreur attendu %v, mais obtenu %v", expectedErrorMessage, receivedErrorMessage)
	}
}

func TestValidateMultipartUpload_TotalSizeTooLarge(t *testing.T) {
	// Simule un multipart upload avec une taille totale dépassant 5 To
	extractPartSizes = func(r *http.Request) []int64 {
		return []int64{2 * 1024 * 1024 * 1024 * 1024, 4 * 1024 * 1024 * 1024 * 1024} // Total = 6 To
	}

	req := httptest.NewRequest(http.MethodPost, "/multipart-upload", nil)

	// Création d'un ResponseRecorder pour capturer la réponse
	rr := httptest.NewRecorder()

	// Utilisation du middleware ValidateMultipartUpload
	handler := ValidateMultipartUpload(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	// Vérification que l'erreur est bien renvoyée
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Code de statut attendu %v, mais obtenu %v", http.StatusRequestEntityTooLarge, rr.Code)
	}

	expectedErrorMessage := ErrMultipartTooLarge.Error()
	receivedErrorMessage := strings.TrimSpace(rr.Body.String())
	if receivedErrorMessage != expectedErrorMessage {
		t.Errorf("Message d'erreur attendu %v, mais obtenu %v", expectedErrorMessage, receivedErrorMessage)
	}
}
