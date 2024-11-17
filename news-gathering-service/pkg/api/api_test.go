package api

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/models"
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/postgres"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_postsHandler(t *testing.T) {
	db, _ := postgres.New("host=localhost user=postgres password=qwerty dbname=news-gathering-service port=5432 sslmode=disable")
	db.CreatePosts([]models.Post{
		{
			Title:   "Title",
			Content: "Content",
			Link:    "Link",
		}})
	api := New(db)

	r := httptest.NewRequest(http.MethodGet, "/news/10", nil)

	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, r)

	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}

	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	var posts []models.Post
	err = json.Unmarshal(b, &posts)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}

	const wantLen = 10
	if len(posts) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(posts), wantLen)
	}
}
