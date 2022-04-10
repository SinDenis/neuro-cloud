package internal

import (
	"demo-rest/internal/config"
	"demo-rest/internal/handler"
	"demo-rest/internal/repository"
	"demo-rest/internal/server"
	"demo-rest/internal/service"
	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	handleErr(container.Provide(config.NewConfig))
	handleErr(container.Provide(config.NewPgPool))
	handleErr(container.Provide(repository.NewUserRepository))
	handleErr(container.Provide(repository.NewImageRepository))
	handleErr(container.Provide(service.NewAuthService))
	handleErr(container.Provide(service.NewRegisterService))
	handleErr(container.Provide(service.NewImageService))
	handleErr(container.Provide(handler.NewAuthHandler))
	handleErr(container.Provide(handler.NewRegisterHandler))
	handleErr(container.Provide(handler.NewImageHandler))
	handleErr(container.Provide(server.NewServer))

	return container
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
