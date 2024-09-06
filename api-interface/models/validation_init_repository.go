package models

import (
	"fmt"
	"reflect"
	"strings"

	entity "api-interface/entities"
	utils "api-interface/utils"
)

// Séquences d'échappement ANSI pour les couleurs
const (
    Reset  = "\033[0m"
    Blue   = "\033[34m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Red    = "\033[31m"
    Grey   = "\033[90m"
)

/** Function retournant les tests d'application d'un TypeEntity avant l'intégration dans la base de données 
    * @param typeName string
    * @param structType reflect.Type
    * @return bool
*/
func AssureEntityCompliance(typeName string, structType reflect.Type) bool {

    fmt.Printf("\n%sTraitement de l'entité : %s%s\n", Blue, typeName, Reset)
    fmt.Println(strings.Repeat("-", 50))

    // Vérifier que le type implémente l'interface Entity
    if !entity.ImplementsEntity(structType) {
        fmt.Printf("%sErreur : L'entité %s n'implémente pas l'interface entity.Entity%s\n", Red, structType.Name(), Reset)
        return false
    }

    // Vérifier que le nom de l'entité est en PascalCase
    if !utils.IsPascalCase(structType.Name()) {
        fmt.Printf("%sErreur : Le nom de l'entité %s n'est pas en PascalCase%s\n", Red, structType.Name(), Reset)
        return false
    }

    // Afficher la structure complète du type
    fmt.Printf("Structure de %s :\n", structType.Name())
    fmt.Println(Grey + strings.Repeat("-", 50) + Reset)
    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        fmt.Printf("  %s%s%s %s%s%s `json:\"%s\"`\n", Blue, field.Name, Reset, Green, field.Type, Reset, field.Tag.Get("json"))
    }

    return true
}