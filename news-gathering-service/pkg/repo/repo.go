package repo

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/models"
)

type DBConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	Port     string `json:"port"`
	SSLMode  string `json:"sslmode"`
}

type Interface interface {
	CreatePosts([]models.Post) error
	Posts(int) ([]models.Post, error)
}
