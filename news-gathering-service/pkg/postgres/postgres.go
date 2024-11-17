package postgres

import (
	"Anastasia/skillfactory/advanced/news-gathering-service/pkg/models"
	"context"
	"math"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

const elemsOnPage = 15

// Конструктор объекта БД
func New(connstr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	s := Store{
		db: db,
	}

	return &s, nil
}

// Загружает в БД посты
func (s *Store) CreatePosts(posts []models.Post) error {
	for _, p := range posts {
		_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (title, content, published_at, link)
		VALUES ($1, $2, $3, $4)`,
			p.Title, p.Content, p.PubTime, p.Link)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Posts(p int, text string) ([]models.Post, models.Pagination, error) {
	var elemsCount int
	row := s.db.QueryRow(context.Background(), `
	SELECT COUNT(id)
	FROM posts
	WHERE title ILIKE $1
	`, "%"+text+"%")

	err := row.Scan(&elemsCount)
	if err != nil {
		return nil, models.Pagination{}, err
	}
	rows, err := s.db.Query(context.Background(), `
	SELECT
		id,
		title,
		published_at,
		link
		FROM posts
		WHERE title ILIKE $1
		ORDER BY published_at DESC
		LIMIT $2
		OFFSET $3`,
		"%"+text+"%", elemsOnPage, (p-1)*elemsOnPage)

	if err != nil {
		return nil, models.Pagination{}, err
	}

	posts := []models.Post{}

	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return nil, models.Pagination{}, err
		}
		posts = append(posts, post)
	}
	pagination := models.Pagination{
		Pages:           int(math.Ceil(float64(elemsCount) / elemsOnPage)),
		CurrentPage:     p,
		ElementsPerPage: elemsOnPage,
	}
	return posts, pagination, rows.Err()
}

func (s *Store) PostDetailed(id int) (models.Post, error) {
	row := s.db.QueryRow(context.Background(), `
		SELECT
		id,
		title,
		content,
		published_at,
		link
		FROM posts
		WHERE id = $1
	`, id)

	post := models.Post{}

	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}
