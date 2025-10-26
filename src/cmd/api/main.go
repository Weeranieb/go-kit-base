package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/weeranieb/go-kit-base/src/internal/config"
	"github.com/weeranieb/go-kit-base/src/internal/di"
	"github.com/weeranieb/go-kit-base/src/internal/handler"
	"github.com/weeranieb/go-kit-base/src/internal/router"

	_ "github.com/weeranieb/go-kit-base/docs" // Import docs package

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

var (
	app *fiber.App
)

var (
	LoadConfigFunc = config.LoadConfig
)

// @title Go Kit Base API
// @version 1.0
// @description A Go application with Fiber, GORM, and Dependency Injection
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

func main() {
	conf := LoadConfigFunc()

	// Dependency Injection
	container := di.NewContainer(conf)

	// Start Fiber + Router
	setupAndStartServer(conf, container)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			log.Println("Gracefully shutting down...")
			shutdownServer()
		}
	}()

	log.Println("Starting server on " + conf.GetServerAddress())
	if err := app.Listen(conf.GetServerAddress()); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func setupAndStartServer(conf *config.Config, container *dig.Container) {
	app = fiber.New(fiber.Config{
		ReadBufferSize: 60 * 1024,
		BodyLimit:      10 * 1024 * 1024, // 10MB
	})

	// Construct the Handler using DI container
	var handlers *handler.Handler

	err := container.Invoke(func(h *handler.Handler) {
		handlers = h
	})
	if err != nil {
		log.Fatal("DI error", err)
	}

	router.SetupRoutes(app, conf, handlers)
	app.Listen(conf.GetServerAddress())
}

func shutdownServer() {
	log.Println("Fiber was successfully shut down.")

	if err := app.Shutdown(); err != nil {
		log.Fatal("Error shutting down Fiber", err)
	}
	os.Exit(0)
}
