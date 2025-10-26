package router

import (
	"github.com/weeranieb/go-kit-base/src/internal/config"
	"github.com/weeranieb/go-kit-base/src/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func NewRouter() *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	return app
}

func SetupRoutes(app *fiber.App, conf *config.Config, handler *handler.Handler) {
	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API routes
	api := app.Group("/api/v1")

	// Setup user routes
	userRouter := NewUserRouter(api)
	userRouter.SetupUserRoutes(handler.UserHandler)
}
