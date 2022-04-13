package main

import (
	"nc-file-processing/internal"
	"nc-file-processing/internal/consumer"
)

func main() {
	container := internal.BuildContainer()
	err := container.Invoke(func(consumer *consumer.ImageConsumer) {
		consumer.Consume()
	})
	if err != nil {
		return
	}
}
