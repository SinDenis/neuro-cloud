package internal

import (
	"go.uber.org/dig"
	"nc-platform-back/internal/config"
	"nc-platform-back/internal/consumer"
	"nc-platform-back/internal/handler"
	"nc-platform-back/internal/producer"
	"nc-platform-back/internal/repository"
	"nc-platform-back/internal/server"
	"nc-platform-back/internal/service"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	handleErr(container.Provide(config.NewConfig))
	handleErr(container.Provide(config.NewPgPool))
	handleErr(container.Provide(config.NewImageClassificationProducer))
	handleErr(container.Provide(config.NewImgClassifyResultConsumer))
	handleErr(container.Provide(repository.NewUserRepository))
	handleErr(container.Provide(repository.NewImageRepository))
	handleErr(container.Provide(service.NewAuthService))
	handleErr(container.Provide(service.NewRegisterService))
	handleErr(container.Provide(service.NewImageService))
	handleErr(container.Provide(handler.NewAuthHandler))
	handleErr(container.Provide(handler.NewRegisterHandler))
	handleErr(container.Provide(handler.NewImageHandler))
	handleErr(container.Provide(server.NewServer))
	handleErr(container.Provide(producer.NewImageProducer))
	handleErr(container.Provide(consumer.NewImageConsumer))

	return container
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
