package middlewares

import (
	"reflect"
	"testing"
)

// TestValidateFilePatternName vérifie la conformité des noms de fichiers selon les règles de nommage.
func TestValidateFilePatternName(t *testing.T) {
	validator := NewFileNameValidator()

	tests := []struct {
		name     string
		expected []string
	}{
		// Cas valides
		{"valid-file-name.txt", []string{"Nom valide."}},
		{"my_image123.png", []string{"Nom valide."}},
		{"Document_final_version.pdf", []string{"Nom valide."}},

		// Cas invalides
		{"", []string{
			"Le nom doit être entre 1 et 1024 caractères.",
			"Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.",
		}}, // Nom vide avec deux erreurs
		{"a", []string{"Nom valide."}}, // Cas valide : un seul caractère est accepté.
		{".hiddenfile", []string{
			"Le nom ne peut pas commencer par un point.",
			"Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.",
		}},
		{"file..name.jpg", []string{"Le nom ne peut pas contenir des points consécutifs."}},
		{"file/with/slash.txt", []string{
			"Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.",
		}},
		{"file-with-emoji-😊.txt", []string{
			"Le nom contient des caractères non valides.",
			"Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.",
		}},
		{"this-file-name-is-way-too-long-" + string(make([]byte, 1000)) + ".txt", []string{
			"Le nom doit être entre 1 et 1024 caractères.",
			"Le nom contient des caractères non valides.",
			"Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.",
		}}, // Nom trop long avec plusieurs erreurs
		{"file\nnewline.txt", []string{
			"Le nom contient des caractères non valides.",
			"Le nom doit contenir uniquement des lettres, chiffres, tirets, underscores, et points.",
		}},
	}

	for _, test := range tests {
		result := validator.Validate(test.name)

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For file name %q, expected errors %v, got %v", test.name, test.expected, result)
		}

		// Affichage des résultats pour tous les cas
		if len(result) == 1 && result[0] == "Nom valide." {
			t.Logf("Is '%s' valid? %v", test.name, result)
		} else {
			t.Logf("For file name '%s', errors: %v", test.name, result)
		}
	}
}
