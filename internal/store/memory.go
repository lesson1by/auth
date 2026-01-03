package store

import (
	"authProject/internal/models"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUsernameRequired  = errors.New("username is required")
)

type UserStore interface {
	Get(username string) (models.User, error)
	Create(user models.User) error
}

type PgxStore struct {
	db *sql.DB
}

func NewPgxStore(db *sql.DB) *PgxStore {
	return &PgxStore{db: db}
}

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		username VARCHAR(255) PRIMARY KEY,
		password VARCHAR(255) NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}

func (p *PgxStore) Get(username string) (models.User, error) {
	var user models.User
	err := p.db.QueryRow(`SELECT username, password FROM users WHERE username = $1 `,
		username,
	).Scan(&user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (p *PgxStore) Create(user models.User) error {
	if user.Username == "" {
		return ErrUsernameRequired
	}
	_, err := p.db.Exec(`INSERT INTO users (username, password) VALUES($1,$2)`,
		user.Username,
		user.Password,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrUserAlreadyExists
			}
		}
		return err
	}
	return nil
}
