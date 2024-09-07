package middlewares

import (
	"api-interface/handlers/errors" // Mettez ici le chemin correct vers votre package d'erreurs
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"net/http"
)

// Exemple de clé secrète. Assurez-vous de gérer les clés secrètes de manière sécurisée.
var (
	secretKey string = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
)

// encrypt chiffre le texte et retourne le texte chiffré ou appelle HandleError en cas d'erreur.
func encrypt(w http.ResponseWriter, plaintext string) (string, bool) {
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
func decrypt(w http.ResponseWriter, ciphertext string) (string, bool) {
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
