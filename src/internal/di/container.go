package di

import (
	"github.com/weeranieb/go-kit-base/src/internal/config"
	"github.com/weeranieb/go-kit-base/src/internal/handler"
	"github.com/weeranieb/go-kit-base/src/internal/repository"
	"github.com/weeranieb/go-kit-base/src/internal/service"

	"go.uber.org/dig"
)

func NewContainer(conf *config.Config) *dig.Container {
	c := dig.New()

	c.Provide(conf.ConnectDB)

	// Repository
	c.Provide(repository.NewUserRepository)

	// Service
	c.Provide(service.NewUserService)

	// Handler
	c.Provide(handler.NewUserHandler)
	c.Provide(handler.NewHandler)

	return c
}
