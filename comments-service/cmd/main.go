package main

import (
	"Anastasia/skillfactory/advanced/comments-service/pkg/api"
	"Anastasia/skillfactory/advanced/comments-service/pkg/postgres"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	api := api.New(db)

	log.Println("Comments service is started on localhost:8082")
	http.ListenAndServe(":8082", api.Router())
}
