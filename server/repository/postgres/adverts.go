package postgres

import (
	"database/sql"
	"fmt"

	"github.com/tousart/marketplace/config"
	"github.com/tousart/marketplace/models"
	pkgDB "github.com/tousart/marketplace/pkg/db"
)

type AdvertsRepo struct {
	db *sql.DB
}

func NewAdvertsRepo(cfg *config.Config) (*AdvertsRepo, error) {
	db, err := pkgDB.ConnectToDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	return &AdvertsRepo{db: db}, nil
}

func (ar *AdvertsRepo) PostAdvert(advert *models.Advert) error {
	query := `INSERT INTO adverts (advert_id, user_id, title, text, url, price, date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := ar.db.Exec(query,
		advert.AdvertID,
		advert.UserID,
		advert.Title,
		advert.Text,
		advert.URL,
		advert.Price,
		advert.Date)
	if err != nil {
		return fmt.Errorf("failed to insert advert to db: %v", err)
	}
	return nil
}
