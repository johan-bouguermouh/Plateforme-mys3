package middlewares

import (
	"api-interface/handlers/errors" // Import du handler d'erreur custom
	"crypto/aes"                    // Package aes implements AES encryption : https://pkg.go.dev/crypto/aes
	"crypto/cipher"                 // Implements standard block cipher modes : https://pkg.go.dev/crypto/cipher
	"crypto/rand"                   // Rand implements a cryptographically secure random number generator : https://pkg.go.dev/crypto/rand

	"github.com/gofiber/fiber/v2" // Fiber is an Express inspired web framework built on top of Fasthttp : https://pkg.go.dev/github.com/gofiber/fiber/v2
)

// Exemple de clé secrète.
var (
	secretKey string = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
)

// encrypt chiffre le texte et retourne le texte chiffré ou appelle HandleError en cas d'erreur.
func encrypt(w *fiber.Ctx, plaintext string) (string, bool) {
	aesBlock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		errors.HandleError(w, errors.ErrInternalServerError, "Encryption failed: "+err.Error())
		return "", false
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		errors.HandleError(w, errors.ErrInternalServerError, "Encryption failed: "+err.Error())
		return "", false
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		errors.HandleError(w, errors.ErrInternalServerError, "Encryption failed: "+err.Error())
		return "", false
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return string(ciphertext), true
}

// decrypt déchiffre le texte et retourne le texte déchiffré ou appelle HandleError en cas d'erreur.
func decrypt(w *fiber.Ctx, ciphertext string) (string, bool) {
	aesBlock, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		errors.HandleError(w, errors.ErrInternalServerError, "Decryption failed: "+err.Error())
		return "", false
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		errors.HandleError(w, errors.ErrInternalServerError, "Decryption failed: "+err.Error())
		return "", false
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		errors.HandleError(w, errors.ErrInternalServerError, "Decryption failed: "+err.Error())
		return "", false
	}

	return string(plaintext), true
}
