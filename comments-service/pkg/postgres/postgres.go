package postgres

import (
	"Anastasia/skillfactory/advanced/comments-service/pkg/models"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

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

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) CreateComment(c models.Comment) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO comments(content, publication_date, news_id)
		VALUES ($1, $2, $3)
	`, c.Content, time.Now().Unix(), c.NewsID)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Comments(newsID int) ([]models.Comment, error) {
	rows, err := s.db.Query(context.Background(), `
        SELECT id, content, publication_date, news_id, moderation_status
        FROM comments
        WHERE news_id = $1
        ORDER BY publication_date ASC
    `, newsID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.Content, &c.PubDate, &c.NewsID, &c.ModStatus)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
