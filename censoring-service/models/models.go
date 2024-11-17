package models

type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	PubDate int64  `json:"publication_date"`
	NewsID  int    `json:"news_id"`
}
