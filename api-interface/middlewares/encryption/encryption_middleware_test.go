package middlewares

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	// Exemple de texte à chiffrer et déchiffrer
	originalText := "This is some sensitive information"

	// Chiffrement
	ciphertext := encrypt(originalText)
	if len(ciphertext) == 0 {
		t.Fatal("Encryption failed: resulting ciphertext is empty")
	}
	t.Logf("Encrypted ciphertext: %x", ciphertext)

	// Déchiffrement
	decryptedText := decrypt(ciphertext)
	if decryptedText != originalText {
		t.Errorf("Decrypted text does not match original text. Got: %s, Want: %s", decryptedText, originalText)
	} else {
		t.Logf("Decrypted text matches original text: %s", decryptedText)
	}

	// Test avec un autre texte
	originalText2 := "Hello"
	ciphertext2 := encrypt(originalText2)
	if len(ciphertext2) == 0 {
		t.Fatal("Encryption failed: resulting ciphertext is empty")
	}
	t.Logf("Encrypted ciphertext 2: %x", ciphertext2)

	decryptedText2 := decrypt(ciphertext2)
	if decryptedText2 != originalText2 {
		t.Errorf("Decrypted text does not match original text. Got: %s, Want: %s", decryptedText2, originalText2)
	} else {
		t.Logf("Decrypted text 2 matches original text: %s", decryptedText2)
	}
}

func TestEncryptDecryptEmpty(t *testing.T) {
	// Test avec un texte vide
	originalText := ""

	// Chiffrement
	ciphertext := encrypt(originalText)
	if len(ciphertext) == 0 {
		t.Fatal("Encryption failed: resulting ciphertext is empty")
	}
	t.Logf("Encrypted empty text ciphertext: %x", ciphertext)

	// Déchiffrement
	decryptedText := decrypt(ciphertext)
	if decryptedText != originalText {
		t.Errorf("Decrypted text does not match original text. Got: %s, Want: %s", decryptedText, originalText)
	} else {
		t.Logf("Decrypted empty text matches original text: %s", decryptedText)
	}
}
