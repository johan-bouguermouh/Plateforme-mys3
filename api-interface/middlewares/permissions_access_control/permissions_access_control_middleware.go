package middlewares

import (
	"api-interface/handlers/errors" // Import du handler d'erreur custom
	"log" // Log (package simple de logging) : https://pkg.go.dev/log
	"net/http" // HTTP client provider : https://pkg.go.dev/net/http
	"strings" // Simple functions to manipulate UTF-8 encoded strings https://pkg.go.dev/strings
)

// Fausse données d'accès (fichier, ou BDD)
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

// PermissionMiddleware 
//  - next (http.Handler)
// Contrôle d'accès sur les chemins de bucket / file
func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// URL - Path 
		path := r.URL.Path
		// Utilisateur en header
		user := r.Header.Get("X-User")
		if user == "" {
			log.Printf("Access denied: Missing user header\n")
			errors.HandleError(w, errors.ErrForbidden, "Missing user header")
			return
		}

		// Extraction du nom du bucket à partir du chemin
		parts := strings.SplitN(path, "/", 3)
		if len(parts) < 2 {
			log.Printf("Access denied: Invalid path format %s\n", path)
			errors.HandleError(w, errors.ErrNotFound, "Invalid path format")
			return
		}

		// Check du bucket
		bucket := parts[1]
		if _, exists := acl[bucket]; !exists {
			// Si le bucket n'existe pas dans l'ACL, le chemin est considéré comme invalide
			log.Printf("Access denied: Bucket %s not found\n", bucket)
			errors.HandleError(w, errors.ErrNotFound) // Utilise le message par défaut
			return
		}

		// Check permission pour l'utilisateur et le bucket
		requiredPermission, userExists := acl[bucket][user]
		if !userExists {
			log.Printf("Access denied: User %s does not have access to bucket %s\n", user, bucket)
			errors.HandleError(w, errors.ErrForbidden, "User does not have access to the bucket")
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

		// Vérification des autorisations en header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer mysecret" {
			log.Printf("Access denied to %s for user %s: Unauthorized\n", path, user)
			errors.HandleError(w, errors.New(http.StatusUnauthorized, "Invalid authorization token"))
			return
		}

		log.Printf("Access granted to %s for user %s\n", path, user)
		next.ServeHTTP(w, r)
	})
}