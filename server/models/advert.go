package models

import "time"

type Advert struct {
	AdvertID int       `json:"advert_id"`
	UserID   string    `json:"user_id"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	URL      string    `json:"url"`
	Price    int       `json:"price"`
	Date     time.Time `json:"date"`
}
