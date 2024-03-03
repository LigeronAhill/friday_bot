package pg

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/LigeronAhill/friday_bot/storage"
	_ "github.com/jackc/pgx"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	log.Print("Connection to db succefull")

	return &Storage{db: db}, nil
}
func (s *Storage) Save(p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES ($1, $2)`

	if _, err := s.db.Exec(q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = $1 ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRow(q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}
func (s *Storage) Remove(page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = $1 AND user_name = $2`
	if _, err := s.db.Exec(q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

func (s *Storage) IsExists(page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = $1 AND user_name = $2`

	var count int

	if err := s.db.QueryRow(q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}
func (s *Storage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name VARCHAR(100))`

	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
