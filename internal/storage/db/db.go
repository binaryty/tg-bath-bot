package db

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

// Article struct.
type Article struct {
	ID          int
	Title       string
	URL         string
	ThumbURL    string
	PublishedAt string
}

// DB struct.
type DB struct {
	db *sql.DB
}

// Id get an id type of string.
func (a Article) Id() string {
	return strconv.Itoa(a.ID)
}

// New ...
func New() (*DB, error) {
	dataSource := "user=postgres password=postgres dbname=articles sslmode=disable"
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

// SaveArticle ...
func (s *DB) SaveArticle(article Article) error {
	query := `INSERT INTO articles(title, url, thumb_url, published_at) VALUES($1, $2, $3, $4) ON CONFLICT DO NOTHING;`
	_, err := s.db.Exec(query, article.Title, article.URL, article.ThumbURL, article.PublishedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetRndArticle ...
func (s *DB) GetRndArticle() (*Article, error) {
	query := `SELECT id, title, url, published_at FROM articles ORDER BY RANDOM() LIMIT 1`

	var art Article

	err := s.db.QueryRow(query).Scan(&art.ID, &art.Title, &art.URL, &art.PublishedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no saved articles: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random article: %v", err)
	}

	return &art, nil
}

// GetArticles ...
func (s *DB) GetArticles() ([]Article, error) {
	query := `SELECT id, title, url FROM articles`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := make([]Article, 0)

	for rows.Next() {
		var art Article
		if err := rows.Scan(&art.ID, &art.Title, &art.URL, &art.PublishedAt); err != nil {
			return nil, err
		}

		articles = append(articles, art)
	}

	return articles, nil
}

// GetArticlesByTitle ...
func (s *DB) GetArticlesByTitle(title string) ([]Article, error) {
	query := `SELECT id, title, url, thumb_url, published_at FROM articles WHERE LOWER(title) LIKE '%' || $1 || '%'	ORDER BY published_at DESC;`

	rows, err := s.db.Query(query, title)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	articles := make([]Article, 0)

	for rows.Next() {
		var art Article
		if err := rows.Scan(&art.ID, &art.Title, &art.URL, &art.ThumbURL, &art.PublishedAt); err != nil {
			return nil, err
		}

		articles = append(articles, art)
	}

	return articles, nil
}
