package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func (ar *AdvertsRepo) GetAdvertsForward(req *models.FeedRequest) ([]models.FeedResponse, error) {
	var (
		lastVal interface{}
		err     error
	)

	if req.FieldSort == "date" {
		lastVal, err = time.Parse("2006-01-02 15:04:05", req.LastVal.(string))
		if err != nil {
			return nil, fmt.Errorf("bad date: %v", err)
		}
	} else if req.FieldSort == "price" {
		lastVal, err = strconv.Atoi(req.LastVal.(string))
		if err != nil {
			return nil, fmt.Errorf("bad price: %v", err)
		}
	}

	var query string

	if strings.Contains(req.TypeSort, "ascending") {
		query = `SELECT advert_id, user_id, title, text, url, price, date FROM adverts 
		WHERE (price BETWEEN $1 AND $2) 
		AND ($3 > $4 OR ($3 = $4 AND advert_id > $5)) 
		ORDER BY $3 ASC, advert_id ASC LIMIT 5;`
	} else if strings.Contains(req.TypeSort, "descending") {
		query = `SELECT advert_id, user_id, title, text, url, price, date FROM adverts 
		WHERE (price BETWEEN $1 AND $2) 
		AND ($3 < $4 OR ($3 = $4 AND advert_id > $5)) 
		ORDER BY $3 DESC, advert_id ASC LIMIT 5;`
	}

	rows, err := ar.db.Query(query,
		req.MinPrice,
		req.MaxPrice,
		req.FieldSort,
		lastVal,
		req.LastID)

	if err != nil {
		return nil, fmt.Errorf("failed to select adverts (forward): %v", err)
	}

	var adverts []models.FeedResponse
	for rows.Next() {
		var ad models.FeedResponse
		err := rows.Scan(&ad.AdvertID, &ad.UserID, &ad.Title, &ad.Text, &ad.URL, &ad.Price, &ad.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row (forward): %v", err)
		}

		if ad.UserID == req.UserID {
			ad.Your = true
		}

		adverts = append(adverts, ad)
	}

	return adverts, nil
}

func (ar *AdvertsRepo) GetAdvertsBackward(req *models.FeedRequest) ([]models.FeedResponse, error) {
	var (
		lastVal interface{}
		err     error
	)

	if req.FieldSort == "date" {
		lastVal, err = time.Parse("2006-01-02 15:04:05", req.LastVal.(string))
		if err != nil {
			return nil, fmt.Errorf("bad date: %v", err)
		}
	} else if req.FieldSort == "price" {
		lastVal, err = strconv.Atoi(req.LastVal.(string))
		if err != nil {
			return nil, fmt.Errorf("bad price: %v", err)
		}
	}

	var query string

	if strings.Contains(req.TypeSort, "ascending") {
		query = `SELECT advert_id, user_id, title, text, url, price, date FROM adverts 
		WHERE (price BETWEEN $1 AND $2) 
		AND ($3 < $4 OR ($3 = $4 AND advert_id < $5)) 
		ORDER BY $3 DESC, advert_id DESC LIMIT 5;`
	} else if strings.Contains(req.TypeSort, "descending") {
		query = `SELECT advert_id, user_id, title, text, url, price, date FROM adverts 
		WHERE (price BETWEEN $1 AND $2) 
		AND ($3 > $4 OR ($3 = $4 AND advert_id < $5)) 
		ORDER BY $3 ASC, advert_id DESC LIMIT 5;`
	}

	rows, err := ar.db.Query(query,
		req.MinPrice,
		req.MaxPrice,
		req.FieldSort,
		lastVal,
		req.LastID)

	if err != nil {
		return nil, fmt.Errorf("failed to select adverts (backward): %v", err)
	}

	var adverts []models.FeedResponse
	for rows.Next() {
		var ad models.FeedResponse
		err := rows.Scan(&ad.AdvertID, &ad.UserID, &ad.Title, &ad.Text, &ad.URL, &ad.Price, &ad.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row (backward): %v", err)
		}

		if ad.UserID == req.UserID {
			ad.Your = true
		}

		adverts = append(adverts, ad)
	}

	return adverts, nil
}
