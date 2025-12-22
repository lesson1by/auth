package store

import (
	"authProject/internal/models"
	"errors"
)

type UserStore interface {
	Get(username string) (models.User, error)
	Create(user models.User) error
}

type InMemoryStore struct {
	Users map[string]models.User
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		Users: make(map[string]models.User),
	}
}

func (in *InMemoryStore) Get(username string) (models.User, error) {
	user, ok := in.Users[username]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (in *InMemoryStore) Create(user models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if _, ok := in.Users[user.Username]; ok {
		return errors.New("user already exists")
	}
	in.Users[user.Username] = user
	return nil
}
