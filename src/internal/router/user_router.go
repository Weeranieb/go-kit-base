package router

import (
	"github.com/weeranieb/go-kit-base/src/internal/handler"

	"github.com/gofiber/fiber/v2"
)

type UserRouter struct {
	group fiber.Router
}

func NewUserRouter(group fiber.Router) *UserRouter {
	return &UserRouter{group: group}
}

func (ur *UserRouter) SetupUserRoutes(userHandler handler.UserHandler) {
	// User routes
	users := ur.group.Group("/users")

	// User CRUD operations
	users.Post("", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
	users.Get("", userHandler.ListUsers)

	users.Get("/:id/profile", userHandler.GetUserProfile)
	users.Put("/:id/profile", userHandler.UpdateUserProfile)
}
