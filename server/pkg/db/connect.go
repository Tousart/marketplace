package pkg

import (
	"database/sql"
	"fmt"

	"github.com/tousart/marketplace/config"
)

func ConnectToDB(cfg *config.Config) (*sql.DB, error) {
	address := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName)

	db, err := sql.Open("postgres", address)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db error: %v", err)
	}

	return db, nil
}
