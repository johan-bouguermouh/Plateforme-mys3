package middlewares

import (
	"regexp"
	"unicode/utf8"
)

// Taille min (1) - max (1024) du nom du fichier.
const (
	minFileNameLength = 1
	maxFileNamelength = 1024
)

// Validator struct encapsule les règles de validation pour les noms.
type Validator struct {
	prefixPattern   *regexp.Regexp
	suffixPattern   *regexp.Regexp
	invalidPattern  *regexp.Regexp
	dotsPattern     *regexp.Regexp
	namePattern     *regexp.Regexp
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

	var errors []string

	// Vérification de la validité UTF-8
	if !utf8.ValidString(name) {
		errors = append(errors, "Le nom doit être encodé en UTF-8.")
	}


	// Vérification de la longueur du nom
	if len(name) < minFileNameLength || len(name) > maxFileNamelength {
		errors = append(errors, "Le nom doit être entre 1 et 1024 caractères.")
	}

	// Vérification pour éviter les noms qui commencent par un point
	if v.prefixPattern.MatchString(name) {
		errors = append(errors, "Le nom ne peut pas commencer par un point.")
	}

	// Vérification des caractères non imprimables ou spéciaux non permis
	if v.invalidPattern.MatchString(name) {
		errors = append(errors, "Le nom contient des caractères non valides.")
	}

	// Vérification pour éviter les points consécutifs
	if v.dotsPattern.MatchString(name) {
		errors = append(errors, "Le nom ne peut pas contenir des points consécutifs.")
	}

	// Vérification du pattern général du nom
	if !v.namePattern.MatchString(name) {
		errors = append(errors, "Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.")
	}

	if len(errors) == 0 {
		return []string{"Nom valide."}
	}

	return errors
}
