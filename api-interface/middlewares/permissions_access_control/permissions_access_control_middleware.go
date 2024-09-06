package middlewares

import (
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extraction du nom du bucket à partir du chemin
		parts := strings.SplitN(path, "/", 3)
		if len(parts) < 2 {
			log.Printf("Access denied: Invalid path format %s\n", path)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		
		// Si le bucket n'existe pas dans l'ACL, le chemin est considéré comme invalide
		bucket := parts[1]
		if _, exists := acl[bucket]; !exists {
			log.Printf("Access denied: Bucket %s not found\n", bucket)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		requiredPermission, userExists := acl[bucket][user]
		if !userExists {
			log.Printf("Access denied: User %s does not have access to bucket %s\n", user, bucket)
			http.Error(w, "Forbidden", http.StatusForbidden)
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
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("Access granted to %s for user %s\n", path, user)
		next.ServeHTTP(w, r)
	})
}
