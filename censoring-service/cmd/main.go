package main

import (
	"Anastasia/skillfactory/advanced/censoring-service/pkg/api"
	"log"
	"net/http"
)

func main() {
	api := api.New()

	log.Println("Censoring service is started on localhost:8083")
	http.ListenAndServe(":8083", api.Router())
}
