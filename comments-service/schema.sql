DROP TABLE IF EXISTS comments;

CREATE TABLE comments(
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    publication_date INTEGER DEFAULT 0,
    news_id INTEGER NOT NULL
);