package storage

import (
	"errors"
	"time"
)

type Storage interface {
	Create(url string) (*Link, error)
	All() []*Link
	Update(slug, url string) error
	Delete(slug string) error
	TryUse(slug string) *Link
}

type Link struct {
	Slug    string    `json:"slug"`
	Url     string    `json:"url"`
	Uses    int       `json:"uses"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

var ErrSlugAlreadyExists = errors.New("slug already exists")
