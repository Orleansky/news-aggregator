package models

type Comment struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	PubDate   int64  `json:"pub_date"`
	NewsID    int    `json:"news_id"`
	ParentID  int    `json:"parent_id"`
	ModStatus string `json:"moderation_status"`
}
