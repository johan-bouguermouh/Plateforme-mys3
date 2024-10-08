package handlers

import (
	"api-interface/database"
	"api-interface/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// UserList returns a list of users
func UserList(c *fiber.Ctx) error {
	// users := database.Get()

	// return c.JSON(fiber.Map{
	// 	"success": true,
	// 	"users":   users,
	// })
	return nil
}

// UserCreate registers a user
func UserCreate(c *fiber.Ctx) error {
	user := &models.User{

		Name: utils.CopyString(c.FormValue("user")),
	}
	database.Insert(user)

	// return c.JSON(fiber.Map{
	// 	"success": true,
	// 	"user":    user,
	// })
	return nil

}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}
