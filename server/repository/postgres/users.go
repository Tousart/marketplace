package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/tousart/marketplace/config"
	pkgDB "github.com/tousart/marketplace/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(cfg *config.Config) (*UsersRepo, error) {
	db, err := pkgDB.ConnectToDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	return &UsersRepo{db: db}, nil
}

func (ur *UsersRepo) Login(login, password string) (string, error) {
	var (
		userID       string
		hashPassword string
	)

	err := ur.db.QueryRow("SELECT user_id, password FROM users WHERE login = $1", login).Scan(&userID, &hashPassword)
	if err == sql.ErrNoRows {
		return "", errors.New("user not exists")
	} else if err != nil {
		return "", fmt.Errorf("select error: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	return userID, nil
}

func (ur *UsersRepo) Register(userID, login, hashPassword string) error {
	_, err := ur.db.Exec("INSERT INTO users (user_id, login, password) VALUES ($1, $2, $3)", userID, login, hashPassword)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errors.New("user exists")
			}
		}
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}
