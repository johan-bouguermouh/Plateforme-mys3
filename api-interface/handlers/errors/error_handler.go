package errors

import (
	"encoding/json" // Librairie encoding JSON : https://pkg.go.dev/encoding/json 
	"net/http" // HTTP client provider : https://pkg.go.dev/net/http
)

// Structure pour une erreur personnalisée : 
// Code : Satus HTTP de l'erreur
// Message : Message d'erreur, par défault il est définit ci-dessous selon le status mais peuvent-être personnalisé
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Nouvelle instance d'erreur
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// HandleError écrit la réponse d'erreur en format JSON dans le writer HTTP.
func HandleError(w http.ResponseWriter, err *Error, customMessage ...string) {

	// Remplacer le message par défaut s'il y a un message personnalisé
	if len(customMessage) > 0 {
		err.Message = customMessage[0]
	}

	// Ajout du status d'erreur aux header de la réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	// Sérialise l'erreur en JSON et l'écrit dans la réponse
	json.NewEncoder(w).Encode(err)
}

// Exemples d'utilisation avec message par défault pour chaque status.
var (
	ErrBadRequest          = New(http.StatusBadRequest, "Bad Request")
	ErrForbidden           = New(http.StatusForbidden, "Forbidden")
	ErrNotFound            = New(http.StatusNotFound, "Not Found")
	ErrInternalServerError = New(http.StatusInternalServerError, "Internal Server Error")
	ErrRequestEntityTooLarge = New(http.StatusRequestEntityTooLarge, "Request entity too large")
)
