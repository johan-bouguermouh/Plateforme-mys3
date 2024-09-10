// handlers/auth.go
package handlers

import (
	"api-interface/database"
	"api-interface/models"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret")

// Inscription d'un utilisateur
func Register(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		log.Printf("Erreur de parsing : %v", err) // Ajoutez ce log
		return c.Status(400).SendString("Données invalides")
	}

	log.Printf("Données reçues : %+v", user)
	// Vérification si l'utilisateur existe déjà
	if database.GetByName(user.Name) != nil {
		return c.Status(400).SendString("L'utilisateur existe déjà")
	}

	// Hachage du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(500).SendString("Erreur lors du hachage du mot de passe")
	}

	// Création de l'utilisateur
	user.Password = string(hashedPassword)
	database.Insert(user)
	return c.SendString("Inscription réussie")
}

// Connexion d'un utilisateur
func Login(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString("Données invalides")
	}

	// Vérification de l'utilisateur
	storedUser := database.GetByName(user.Name)
	if storedUser == nil {
		return c.Status(400).SendString("Utilisateur ou mot de passe incorrect")
	}

	// Vérification du mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return c.Status(400).SendString("Utilisateur ou mot de passe incorrect")
	}

	// Génération d'un token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Name,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).SendString("Erreur lors de la génération du token")
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		log.Println("Missing Authorization Header")
		return c.Status(fiber.StatusUnauthorized).SendString("Missing Authorization Header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		log.Println("Missing Bearer Token")
		return c.Status(fiber.StatusUnauthorized).SendString("Missing Bearer Token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("Token parsing error:", err)
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
	}

	if !token.Valid {
		log.Println("Invalid token")
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
	}

	return c.Next()
}
