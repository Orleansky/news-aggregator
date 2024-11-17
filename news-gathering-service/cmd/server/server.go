package main

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/api"
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/models"
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/postgres"
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/repo"
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/rss"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	RSS    []string      `json:"rss"`
	Period int           `json:"request_period"`
	DB     repo.DBConfig `json:"db"`
}

func main() {
	b, err := os.ReadFile("C:/Users/LK10725/go/src/Anastasia/skillfactory/advanced/news-gathering-service/cmd/server/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var c config
	err = json.Unmarshal(b, &c)

	if err != nil {
		log.Fatal(err.Error())
	}

	connstr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DB.Host, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.Port, c.DB.SSLMode)

	db, err := postgres.New(connstr)
	if err != nil {
		log.Fatal(err.Error())
	}

	api := api.New(db)

	chanPosts := make(chan []models.Post)
	chanErrors := make(chan error)

	for _, url := range c.RSS {
		go parse(url, chanPosts, chanErrors, c.Period)
	}

	go func() {
		for posts := range chanPosts {
			db.CreatePosts(posts)
		}
	}()

	go func() {
		for err := range chanErrors {
			log.Print(err.Error())
		}
	}()

	log.Println("News gathering service is started on localhost:8081")
	http.ListenAndServe(":8081", api.Router())
}

func parse(url string, posts chan<- []models.Post, errs chan<- error, period int) {

	for range time.Tick(time.Duration(period) * time.Second) {
		news, err := rss.Read(url)
		if err != nil {
			errs <- err
			continue
		}

		posts <- news
	}
}
