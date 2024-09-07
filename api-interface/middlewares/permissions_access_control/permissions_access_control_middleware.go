package middlewares

import (
	"api-interface/handlers/errors"
	"log"
	"net/http"
	"strings"
)

// Fausse base de données des permissions
var acl = map[string]map[string]string{
	"bucket1": {
		"user1": "public-read",
		"user2": "private",
	},
	"bucket2": {
		"user3": "public-read",
		"user4": "private",
	},
}

// PermissionMiddleware contrôle d'accès sur les chemins de bucket / file
func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		user := r.Header.Get("X-User")
		if user == "" {
			log.Printf("Access denied: Missing user header\n")
			errors.HandleError(w, errors.ErrForbidden, "Missing user header") // Message personnalisé
			return
		}

		// Extraction du nom du bucket à partir du chemin
		parts := strings.SplitN(path, "/", 3)
		if len(parts) < 2 {
			log.Printf("Access denied: Invalid path format %s\n", path)
			errors.HandleError(w, errors.ErrNotFound, "Invalid path format") // Message personnalisé
			return
		}

		bucket := parts[1]
		if _, exists := acl[bucket]; !exists {
			// Si le bucket n'existe pas dans l'ACL, le chemin est considéré comme invalide
			log.Printf("Access denied: Bucket %s not found\n", bucket)
			errors.HandleError(w, errors.ErrNotFound) // Utilise le message par défaut
			return
		}

		requiredPermission, userExists := acl[bucket][user]
		if !userExists {
			log.Printf("Access denied: User %s does not have access to bucket %s\n", user, bucket)
			errors.HandleError(w, errors.ErrForbidden, "User does not have access to the bucket") // Message personnalisé
			return
		}

		// Log de la permission requise
		log.Printf("User %s requires %s permission for bucket %s\n", user, requiredPermission, bucket)

		// Vérifiez la permission
		if requiredPermission == "public-read" {
			log.Printf("Access granted to %s for user %s\n", path, user)
			next.ServeHTTP(w, r)
			return
		}

		// Vérification des autorisations
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer mysecret" {
			log.Printf("Access denied to %s for user %s: Unauthorized\n", path, user)
			errors.HandleError(w, errors.New(http.StatusUnauthorized, "Invalid authorization token")) // Utilise 401 Unauthorized
			return
		}

		log.Printf("Access granted to %s for user %s\n", path, user)
		next.ServeHTTP(w, r)
	})
}