package main

import (
	"demo-rest/internal"
	"demo-rest/internal/server"
)

func main() {
	container := internal.BuildContainer()

	err := container.Invoke(func(server *server.Server) {
		server.Run()
	})
	if err != nil {
		panic(err)
	}
}
