package middlewares

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	// Taille maximale des fichiers upload direct -> 5 Go.
	MaxDirectUploadSize = 5 * 1024 * 1024 * 1024
	// Taille maximale avec Multipart Upload : 5 To.
	MaxMultipartUploadSize = 5 * 1024 * 1024 * 1024 * 1024
	// Taille minimale pour les fichiers en Multipart Upload ( sauf dernier qui peut-être plus petit)
	MinPartSize = 5 * 1024 * 1024 // 5 Mo
)

// Error messages
var (
	ErrFileTooLarge      = errors.New("le fichier dépasse la taille maximale autorisée pour un upload direct (5 Go)")
	ErrMultipartTooLarge = errors.New("la taille cumulée des parties dépasse la limite de 5 To pour un Multipart Upload")
	ErrPartTooSmall      = errors.New("chaque partie de l'upload Multipart doit faire au moins 5 Mo, sauf la dernière partie")
)

// ValidateDirectUpload vérifie la taille du fichier pour un upload en "direct" 
func ValidateDirectUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxDirectUploadSize)
		contentLengthStr := r.Header.Get("Content-length")

		if contentLengthStr != "" {
			contentLengthStr, err := strconv.ParseInt(contentLengthStr, 10, 64)
			if err != nil {
				http.Error(w, "Taille du fichier invalide", http.StatusBadRequest)
				return
			}

			if contentLengthStr > MaxDirectUploadSize {
				http.Error(w, ErrFileTooLarge.Error(), http.StatusRequestEntityTooLarge)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateMultipartUpload vérifie la taille d'un fichier upload en multipart 
func ValidateMultipartUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simuler l'extraction des tailles des parties
		partSizes := extractPartSizes(r) // Cette fonction devrait être implémentée pour extraire les tailles des parties d'un upload multipart.

		var totalSize int64

		for i, partSize := range partSizes {
			totalSize += partSize

			// Vérifier la taille de chaque partie du fichier.
			if i < len(partSizes)-1 && partSize < MinPartSize {
				http.Error(w, ErrPartTooSmall.Error(), http.StatusRequestEntityTooLarge)
				return
			}

		}

		// Vérifiez que la taille du fichier ne dépasse pas 5 To
		if totalSize > MaxMultipartUploadSize {
			http.Error(w, ErrMultipartTooLarge.Error(), http.StatusRequestEntityTooLarge)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Simuler l'extraction des tailles des parties d'un upload multipart
// Cette fonction est un exemple, et vous devrez l'adapter à votre logique spécifique
var extractPartSizes = func(r *http.Request) []int64 {
	return []int64{10 * 1024 * 1024, 6 * 1024 * 1024, 3 * 1024 * 1024} // Simule 3 parties de 10 Mo, 6 Mo, et 3 Mo
}
