package models

type NewsFullDetailed struct {
	ID      int
	Title   string
	Content string
	PubDate int64
	Link    string
}

type NewsShortDetailed struct {
	ID      int
	Title   string
	PubDate int64
	Link    string
}

type Comment struct {
	ID      int
	Content string
	PubDate int64
	NewsID  int
}
