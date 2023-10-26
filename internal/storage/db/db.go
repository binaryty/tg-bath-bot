package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Article struct {
	ID    int
	Title string
	URL   string
}

type DB struct {
	db *sql.DB
}

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

func (s *DB) SaveArticle(article Article) error {
	query := `INSERT INTO articles(title, url) VALUES($1, $2)`
	_, err := s.db.Exec(query, article.Title, article.URL)
	if err != nil {
		return err
	}

	return nil
}

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
		if err := rows.Scan(&art.ID, &art.Title, &art.URL); err != nil {
			return nil, err
		}

		articles = append(articles, art)
	}

	return articles, nil
}
