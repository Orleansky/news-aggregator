package rss

import (
	"testing"
)

func TestRead(t *testing.T) {
	posts, err := Read("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		t.Error(err)
	}

	if len(posts) == 0 {
		t.Error("Данные не получены")
	}

	t.Logf("Получено %d статей", len(posts))
}
