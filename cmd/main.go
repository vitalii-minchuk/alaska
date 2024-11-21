package main

import (
	"log"

	"github.com/vitalii-minchuk/alaska/cmd/api"
)

func main() {
	apiServer := api.NewAPIServer(":8080")
	if err := apiServer.Run(); err != nil {
		log.Fatal("error running api server")
	}
}
