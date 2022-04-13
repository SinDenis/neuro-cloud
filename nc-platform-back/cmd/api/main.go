package main

import (
	"fmt"
	"nc-platform-back/internal"
	"nc-platform-back/internal/consumer"
	"nc-platform-back/internal/server"
)

func main() {
	container := internal.BuildContainer()

	go func() {
		err := container.Invoke(func(consumer *consumer.ImageClassResultConsumer) {
			consumer.Consume()
		})
		if err != nil {
			fmt.Println("Consumer dead")
		}
	}()
	err := container.Invoke(func(server *server.Server) {
		server.Run()
	})
	if err != nil {
		panic(err)
	}
}
