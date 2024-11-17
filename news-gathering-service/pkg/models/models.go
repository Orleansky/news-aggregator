package models

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	PubTime int64  `json:"pub_time"`
	Link    string `json:"link"`
}

type Pagination struct {
	Pages           int
	CurrentPage     int
	ElementsPerPage int
}
