package internal

import (
	"go.uber.org/dig"
	"nc-file-processing/internal/config"
	"nc-file-processing/internal/consumer"
	"nc-file-processing/internal/producer"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	handleErr(container.Provide(config.NewConfig))
	handleErr(container.Provide(config.NewImageClassificationKafkaConsumer))
	handleErr(container.Provide(config.NewClassifyImgResultProducer))
	handleErr(container.Provide(consumer.NewImageConsumer))
	handleErr(container.Provide(producer.NewImageProducer))

	return container
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
