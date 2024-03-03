package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/LigeronAhill/friday_bot/lib/e"
)

var (
	ErrNoSavedPages = errors.New("no saved page")
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (res string, err error) {
	defer func() { err = e.Wrap("can't hash", err) }()
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", err
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
