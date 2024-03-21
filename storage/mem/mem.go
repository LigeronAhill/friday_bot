package mem

import (
	"math/rand"

	"github.com/LigeronAhill/friday_bot/storage"
)

type Storage struct {
	db []*storage.Page
}

func New() *Storage {
	db := make([]*storage.Page, 0)
	return &Storage{db: db}
}

func (s *Storage) Save(p *storage.Page) error {

	s.db = append(s.db, p)

	return nil
}

func (s *Storage) PickRandom(userName string) (*storage.Page, error) {
	usersStorage := make([]*storage.Page, 0)
	for _, p := range s.db {
		if p.UserName == userName {
			usersStorage = append(usersStorage, p)
		}
	}
	if len(usersStorage) > 0 {
		index := rand.Intn(len(usersStorage))
		return usersStorage[index], nil
	} else {
		return nil, storage.ErrNoSavedPages
	}
}
func (s *Storage) Remove(page *storage.Page) error {
	var index = 0
	for i, p := range s.db {
		if p == page {
			index = i
			break
		}
	}
	s.db = append(s.db[:index], s.db[index+1:]...)
	return nil
}

func (s *Storage) IsExists(page *storage.Page) (bool, error) {
	for _, p := range s.db {
		if p == page {
			return true, nil
		}
	}
	return false, nil
}
