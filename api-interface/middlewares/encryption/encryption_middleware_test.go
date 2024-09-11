package middlewares

import (
	"net/http/httptest"
	"testing"
)

// Helper function to simulate HTTP response writer
func newTestResponseWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestEncryptDecrypt(t *testing.T) {
	// Exemple de texte à chiffrer et déchiffrer
	originalText := "This is some sensitive information"

	// Créer un ResponseRecorder pour capturer les erreurs
	recorder := newTestResponseWriter()

	// Chiffrement
	ciphertext, success := encrypt(recorder, originalText)
	if !success {
		t.Fatalf("Encryption failed with status: %d", recorder.Code)
	}
	if len(ciphertext) == 0 {
		t.Fatal("Encryption failed: resulting ciphertext is empty")
	}
	t.Logf("Encrypted ciphertext: %x", ciphertext)

	// Déchiffrement
	recorder = newTestResponseWriter() // Reset recorder for the next operation
	decryptedText, success := decrypt(recorder, ciphertext)
	if !success {
		t.Fatalf("Decryption failed with status: %d", recorder.Code)
	}
	if decryptedText != originalText {
		t.Errorf("Decrypted text does not match original text. Got: %s, Want: %s", decryptedText, originalText)
	} else {
		t.Logf("Decrypted text matches original text: %s", decryptedText)
	}

	// Test avec un autre texte
	originalText2 := "Hello"
	recorder = newTestResponseWriter() // Reset recorder for the next operation
	ciphertext2, success := encrypt(recorder, originalText2)
	if !success {
		t.Fatalf("Encryption failed with status: %d", recorder.Code)
	}
	if len(ciphertext2) == 0 {
		t.Fatal("Encryption failed: resulting ciphertext is empty")
	}
	t.Logf("Encrypted ciphertext 2: %x", ciphertext2)

	recorder = newTestResponseWriter() // Reset recorder for the next operation
	decryptedText2, success := decrypt(recorder, ciphertext2)
	if !success {
		t.Fatalf("Decryption failed with status: %d", recorder.Code)
	}
	if decryptedText2 != originalText2 {
		t.Errorf("Decrypted text does not match original text. Got: %s, Want: %s", decryptedText2, originalText2)
	} else {
		t.Logf("Decrypted text 2 matches original text: %s", decryptedText2)
	}
}

func TestEncryptDecryptEmpty(t *testing.T) {
	// Test avec un texte vide
	originalText := ""

	// Créer un ResponseRecorder pour capturer les erreurs
	recorder := newTestResponseWriter()

	// Chiffrement
	ciphertext, success := encrypt(recorder, originalText)
	if !success {
		t.Fatalf("Encryption failed with status: %d", recorder.Code)
	}
	if len(ciphertext) == 0 {
		t.Fatal("Encryption failed: resulting ciphertext is empty")
	}
	t.Logf("Encrypted empty text ciphertext: %x", ciphertext)

	// Déchiffrement
	recorder = newTestResponseWriter() // Reset recorder for the next operation
	decryptedText, success := decrypt(recorder, ciphertext)
	if !success {
		t.Fatalf("Decryption failed with status: %d", recorder.Code)
	}
	if decryptedText != originalText {
		t.Errorf("Decrypted text does not match original text. Got: %s, Want: %s", decryptedText, originalText)
	} else {
		t.Logf("Decrypted empty text matches original text: %s", decryptedText)
	}
}
