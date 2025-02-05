package storage

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Mysql struct {
	db *sqlx.DB

	slugLength int
}

func NewMysql(connect string, slugLength int) (*Mysql, error) {
	dsn := fmt.Sprintf("%s?parseTime=true", connect)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	schema := `CREATE TABLE IF NOT EXISTS links (
        slug VARCHAR(255) NOT NULL,
        url TEXT NOT NULL,
        uses INT NOT NULL DEFAULT 0,
        created DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (slug)
    )`
	if _, err = db.Exec(schema); err != nil {
		return nil, err
	}

	return &Mysql{db: db, slugLength: slugLength}, nil
}

func (s *Mysql) Create(url string) (*Link, error) {
	slug := generateSlug(s.slugLength)

	_, err := s.db.Exec("INSERT INTO links (slug, url) VALUES (?, ?)", slug, url)
	if err != nil {
		if isDuplicate(err) {
			return nil, ErrSlugAlreadyExists
		}
		return nil, err
	}

	return &Link{
		Slug:    slug,
		Url:     url,
		Uses:    0,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil
}

func (s *Mysql) All() []*Link {
	links := make([]*Link, 0)
	if err := s.db.Select(&links, "SELECT * FROM links"); err != nil {
		return links
	}
	return links
}

func (s *Mysql) Update(slug, url string) error {
	_, err := s.db.Exec("UPDATE links SET url = ? WHERE slug = ?", url, slug)
	if err != nil {
		log.Printf("mysql err: %v", err)
		return err
	}
	return nil
}

func (s *Mysql) Delete(slug string) error {
	_, err := s.db.Exec("DELETE FROM links WHERE slug = ?", slug)
	if err != nil {
		log.Printf("mysql err: %v", err)
		return err
	}
	return nil
}

func (s *Mysql) TryUse(slug string) *Link {
	var link Link
	if err := s.db.Get(&link, "SELECT * FROM links WHERE slug = ?", slug); err != nil {
		log.Printf("mysql err: %v", err)
		return nil // no link found
	}
	link.Uses++

	_, err := s.db.Exec("UPDATE links SET uses = ? WHERE slug = ?", link.Uses, slug)
	if err != nil {
		log.Printf("mysql error: %v", err)
		return nil
	}

	return &link
}

func (s *Mysql) Close() error {
	return s.db.Close()
}

func isDuplicate(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

// https://stackoverflow.com/a/31832326
func generateSlug(n int) string {
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
