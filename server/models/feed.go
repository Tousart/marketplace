package models

import "time"

type FeedRequest struct {
	TypeSort  string      `json:"type_sort"`
	FieldSort string      `json:"sort_field"`
	MinPrice  int         `json:"min_price"`
	MaxPrice  int         `json:"max_price"`
	LastVal   interface{} `json:"last_value"`
	LastID    int         `json:"last_id"`
	UserID    string      `json:"user_id"`
}

type FeedResponse struct {
	AdvertID int       `json:"advert_id"`
	UserID   string    `json:"user_id"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	URL      string    `json:"url"`
	Price    int       `json:"price"`
	Date     time.Time `json:"date"`
	Your     bool      `json:"your"`
}
