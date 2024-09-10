package middlewares

import (
	"api-interface/handlers/errors" // Import du handler d'erreur custom
	"net/http" // HTTP client provider : https://pkg.go.dev/net/http
	"regexp" // Regular expression search : https://pkg.go.dev/regexp
	"strings" // Simple functions to manipulate UTF-8 encoded strings https://pkg.go.dev/strings
	"unicode/utf8" // Functions and constants to support text encoded in UTF-8 : https://pkg.go.dev/unicode/utf8
)

// Taille min (1) - max (1024) du nom du fichier.
const (
	minFileNameLength = 1
	maxFileNameLength = 1024
)

// Validator struct encapsule les règles de validation pour les noms.
type Validator struct {
	prefixPattern  *regexp.Regexp
	suffixPattern  *regexp.Regexp
	invalidPattern *regexp.Regexp
	dotsPattern    *regexp.Regexp
	namePattern    *regexp.Regexp
}

// NewFileNameValidator crée une nouvelle instance de Validator avec des règles prédéfinies.
func NewFileNameValidator() *Validator {
	return &Validator{
		// Aucun préfixe spécifique n'est interdit. On évite que le nom commence par un point.
		prefixPattern: regexp.MustCompile(`^\.`),

		// Pas de suffixe spécifique interdit, donc on ne modifie pas.
		suffixPattern: regexp.MustCompile(`$`),

		// Vérifie les caractères non imprimables ou caractères spéciaux non permis dans les systèmes de fichiers
		invalidPattern: regexp.MustCompile(`[^\p{L}\p{N}\p{P}\p{Z}]`),

		// Vérifie les points consécutifs, qui ne sont pas autorisés.
		dotsPattern: regexp.MustCompile(`\.\.`),

		// Caractères autorisés : Lettres minuscules/majuscules, chiffres, tirets, underscores et points
		namePattern: regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`), // Autorise les lettres, chiffres, tirets, underscores, points
	}
}

// Validate vérifie la conformité du nom selon les règles de nommage.
func (v *Validator) Validate(name string) []string {

	var validationErrors []string

	// Vérification de la validité UTF-8
	if !utf8.ValidString(name) {
		validationErrors = append(validationErrors, "Le nom doit être encodé en UTF-8.")
	}

	// Vérification de la longueur du nom
	if len(name) < minFileNameLength || len(name) > maxFileNameLength {
		validationErrors = append(validationErrors, "Le nom doit être entre 1 et 1024 caractères.")
	}

	// Vérification pour éviter les noms qui commencent par un point
	if v.prefixPattern.MatchString(name) {
		validationErrors = append(validationErrors, "Le nom ne peut pas commencer par un point.")
	}

	// Vérification des caractères non imprimables ou spéciaux non permis
	if v.invalidPattern.MatchString(name) {
		validationErrors = append(validationErrors, "Le nom contient des caractères non valides.")
	}

	// Vérification pour éviter les points consécutifs
	if v.dotsPattern.MatchString(name) {
		validationErrors = append(validationErrors, "Le nom ne peut pas contenir des points consécutifs.")
	}

	// Vérification du pattern général du nom
	if !v.namePattern.MatchString(name) {
		validationErrors = append(validationErrors, "Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.")
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return []string{"Nom valide."}
}

// FileNameValidationMiddleware vérifie la validité du nom de fichier dans la requête.
func FileNameValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Supposons que le nom de fichier est passé en tant que paramètre de requête "filename"
		fileName := r.URL.Query().Get("filename")

		// Crée un nouvel instance du validateur
		validator := NewFileNameValidator()

		// Valide le nom du fichier
		nameError := validator.Validate(fileName)

		if len(nameError) > 1 {
			// Si il y a des erreurs, retourne une réponse d'erreur
			errors.HandleError(w, errors.ErrBadRequest, strings.Join(nameError, ", ")) // Utilise une erreur 400 Bad Request
			return
		}

		// Continue avec le traitement suivant si le nom de fichier est valide
		next.ServeHTTP(w, r)
	})
}
