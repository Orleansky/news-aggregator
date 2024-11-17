package rss

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/models"
	"encoding/xml"
	"io"
	"net/http"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

type Sourse struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title   string `xml:"title"`
	Content string `xml:"description"`
	Link    string `xml:"link"`
	Items   []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"description"`
	PubTime string `xml:"pubDate"`
	Link    string `xml:"link"`
}

// Чтение и десериализация данных
// формата XML из RSS-потока
func Read(url string) ([]models.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var s Sourse
	xml.Unmarshal(b, &s)

	posts := []models.Post{}

	for _, item := range s.Channel.Items {
		var p models.Post
		p.Title = item.Title
		p.Content = strip.StripTags(item.Content)
		p.Link = item.Link

		t, err := time.Parse(time.RFC1123Z, item.PubTime)
		if err != nil {
			t, err = time.Parse(time.RFC1123, item.PubTime)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}

		posts = append(posts, p)
	}
	return posts, nil
}
