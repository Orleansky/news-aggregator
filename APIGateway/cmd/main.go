package main

import (
	"Anastasia/skillfactory/advanced/APIGateway/pkg/api"
	"log"
	"net/http"
)

func main() {
	api := api.New()

	log.Println("API Gateway is started on localhost:8080")
	http.ListenAndServe(":8080", api.Router())
}
