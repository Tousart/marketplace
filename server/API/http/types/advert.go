package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AdvertRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
	Price string `json:"price"`
}

type AdvertResponse struct {
	AdvertID string `json:"advert_id"`
	UserID   string `json:"user_id"`
	Title    string `json:"title"`
	Date     string `json:"date"`
}

type AdvertsFeedResponse struct {
	TypeSort  string `json:"type_sort"`
	FieldSort string `json:"sort_field"`
	MinPrice  string `json:"min_price"`
	MaxPrice  string `json:"max_price"`
	LastVal   string `json:"last_value"`
	LastID    string `json:"last_id"`
}

func CreateAdvertRequestHandler(r *http.Request) (*AdvertRequest, error) {
	var advert AdvertRequest
	err := json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode advert request: %v", err)
	}
	return &advert, nil
}

func CreateAdvertsFeedHandler(r *http.Request) (*AdvertsFeedResponse, error) {
	typeSort := r.URL.Query().Get("type_sort")
	fieldSort := r.URL.Query().Get("field_sort")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	lastVal := r.URL.Query().Get("last_value")
	lastID := r.URL.Query().Get("last_id")

	return &AdvertsFeedResponse{
		TypeSort:  typeSort,
		FieldSort: fieldSort,
		MinPrice:  minPriceStr,
		MaxPrice:  maxPriceStr,
		LastVal:   lastVal,
		LastID:    lastID,
	}, nil
}
