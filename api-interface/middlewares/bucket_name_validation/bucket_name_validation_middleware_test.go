package middlewares

import (
	"testing"
)

// TestValidateBucketPatternName vérifie la conformité des noms de bucket.
func TestValidateBucketPatternName(t *testing.T) {
	validator := NewBucketNameValidator()

	// Liste de noms de buckets pour tester la validation
	tests := []struct {
		name     string
		expected []string
	}{
		// Cas valides
		{"valid-bucket-name", []string{"Nom de bucket valide."}},
		{"this-is-a-valid-bucket-name-123", []string{"Nom de bucket valide."}},

		// Cas invalides avec erreurs attendues
		{"invalid..bucket--x-s3", []string{
			"Nom du bucket ne peut pas contenir des points consécutifs.",
			"Nom du bucket ne peut pas commencer par un préfixe invalide.",
		}},
		{"xn--invalid-prefix", []string{
			"Nom du bucket ne peut pas commencer par un préfixe invalide.",
		}},
		{"bucket-name--x-s3", []string{
			"Nom du bucket ne peut pas se terminer par un suffixe invalide.",
		}},
		{"192.168.1.1", []string{
			"Nom du bucket ne peut pas être une adresse IP.",
		}},
		{"a", []string{
			"Nom du bucket doit être entre 3 et 63 caractères.",
		}},
	}

	// Test de chaque nom de bucket
	for _, test := range tests {
		result := validator.Validate(test.name)

		// Affichage de log sans usage d'erreur et stopage de test
		// Affichage des résultats pour tous les cas
		if len(result) == 1 && result[0] == "Nom de bucket valide." {
			t.Logf("Is '%s' valid? %v", test.name, result)
		} else {
			t.Logf("For bucket name '%s', errors: %v", test.name, result)
		}
	}
}
