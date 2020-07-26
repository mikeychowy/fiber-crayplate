package routes

import (
	Controller "github.com/mikeychowy/fiber-crayplate/app/controllers/api"

	"github.com/gofiber/fiber"
)

// RegisterAPI Register All API Routes.
func RegisterAPI(api fiber.Router) {
	registerUsers(api)
}

func registerUsers(api fiber.Router) {
	users := api.Group("/users")

	users.Get("/", Controller.GetAllUsers)
	users.Get("/:id", Controller.GetUser)
	users.Post("/", Controller.AddUser)
	users.Put("/:id", Controller.EditUser)
	users.Delete("/:id", Controller.DeleteUser)
}
