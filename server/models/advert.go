package models

import "time"

type Advert struct {
	AdvertID string
	UserID   string
	Title    string
	Text     string
	URL      string
	Price    int
	Date     time.Time
}
