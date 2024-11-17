package postgres

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/models"
	"math/rand"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	connstr := "host=localhost user=postgres password=qwerty dbname=news-gathering-service port=5432 sslmode=disable"
	_, err := New(connstr)
	if err != nil {
		t.Error(err)
	}
}

func TestStore_CreatePosts(t *testing.T) {
	var posts []models.Post
	for i := 0; i < rand.Intn(100); i++ {
		posts = append(posts, models.Post{
			Title:   "Title",
			Content: "Content",
			PubTime: rand.Int63n(100000000),
			Link:    strconv.Itoa(rand.Intn(100000000)),
		})
	}

	db, err := New("host=localhost user=postgres password=qwerty dbname=news-gathering-service port=5432 sslmode=disable")
	if err != nil {
		t.Error(err)
	}

	err = db.CreatePosts(posts)
	if err != nil {
		t.Error(err)
	}
}

func TestStore_Posts(t *testing.T) {
	db, err := New("host=localhost user=postgres password=qwerty dbname=news-gathering-service port=5432 sslmode=disable")
	if err != nil {
		t.Error(err)
	}

	posts, pagination, err := db.Posts(2, "")
	if err != nil {
		t.Error(err)
	}

	t.Logf("Получено %d постов", len(posts))
	t.Logf("Pagination: %v", pagination)
}
