package middlewares

import (
	"api-interface/handlers/errors" // Import du handler d'erreur custom
	"net/http" // HTTP client provider : https://pkg.go.dev/net/http
	"strconv" // Conversions to and from string representations of basic data types : https://pkg.go.dev/strconv
)

// Const de taille de fichiers.
const (
	// Taille maximale des fichiers upload direct -> 5 Go.
	MaxDirectUploadSize = 5 * 1024 * 1024 * 1024

	// Taille maximale avec Multipart Upload : 5 To.
	MaxMultipartUploadSize = 5 * 1024 * 1024 * 1024 * 1024

	// Taille minimale pour les fichiers en Multipart Upload ( sauf dernier qui peut-être plus petit)
	MinPartSize = 5 * 1024 * 1024 // 5 Mo
)

// ValidateDirectUpload vérifie la taille du fichier pour un upload en "direct".
func ValidateDirectUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxDirectUploadSize)
		contentLengthStr := r.Header.Get("Content-length")

		if contentLengthStr != "" {
			contentLengthStr, err := strconv.ParseInt(contentLengthStr, 10, 64)
			if err != nil {
				errors.HandleError(w, errors.ErrBadRequest, "Taille du fichier invalide")
				return
			}

			// Fichier trop gros
			if contentLengthStr > MaxDirectUploadSize {
				errors.HandleError(w, errors.ErrRequestEntityTooLarge, "Le fichier dépasse la taille maximale autorisée pour un upload direct (5 Go)")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// ValidateMultipartUpload vérifie la taille d'un fichier upload en multipart.
func ValidateMultipartUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Simuler l'extraction des tailles des parties du fichier
		partSizes := extractPartSizes(r)

		var totalSize int64

		for i, partSize := range partSizes {
			totalSize += partSize

			// Vérifier la taille de chaque partie du fichier.
			if i < len(partSizes)-1 && partSize < MinPartSize {
				errors.HandleError(w, errors.ErrRequestEntityTooLarge, "Chaque partie de l'upload Multipart doit faire au moins 5 Mo, sauf la dernière partie")
				return
			}

		}

		// Vérifiez que la taille du fichier ne dépasse pas 5 To
		if totalSize > MaxMultipartUploadSize {
			errors.HandleError(w, errors.ErrRequestEntityTooLarge, "La taille cumulée des parties dépasse la limite de 5 To pour un Multipart Upload")
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
