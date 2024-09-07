package errors

import (
	"encoding/json"
	"net/http"
)

// Error représente une erreur personnalisée avec un code de statut HTTP.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// New crée une nouvelle instance d'erreur avec un message par défaut.
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// HandleError écrit la réponse d'erreur en format JSON dans le writer HTTP.
// Permet de passer un message personnalisé qui remplace le message par défaut si spécifié.
func HandleError(w http.ResponseWriter, err *Error, customMessage ...string) {
	// Remplacer le message par défaut s'il y a un message personnalisé
	if len(customMessage) > 0 {
		err.Message = customMessage[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	// Sérialise l'erreur en JSON et l'écrit dans la réponse
	json.NewEncoder(w).Encode(err)
}

// Exemples d'instances d'erreurs avec des messages par défaut
var (
	ErrBadRequest          = New(http.StatusBadRequest, "Bad Request")
	ErrForbidden           = New(http.StatusForbidden, "Forbidden")
	ErrNotFound            = New(http.StatusNotFound, "Not Found")
	ErrInternalServerError = New(http.StatusInternalServerError, "Internal Server Error")
	ErrRequestEntityTooLarge = New(http.StatusRequestEntityTooLarge, "Request entity too large")
)
