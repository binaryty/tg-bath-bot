package db

import (
	"database/sql"
	"fmt"
	"strconv"

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

func (a Article) Id() string {
	return strconv.Itoa(a.ID)
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

func (s *DB) GetRndArticle() (*Article, error) {
	query := `SELECT id, title, url FROM articles ORDER BY RANDOM() LIMIT 1`

	var art Article

	err := s.db.QueryRow(query).Scan(&art.ID, &art.Title, &art.URL)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no saved articles: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random article: %v", err)
	}

	return &art, nil
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

func (s *DB) GetArticlesByTitle(title string, limit int) ([]Article, error) {
	query := `SELECT id, title, url FROM articles WHERE LOWER(title) LIKE '%' || $1 || '%'	ORDER BY title DESC LIMIT $2;`

	rows, err := s.db.Query(query, title, limit)
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
