package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateAdvertRequestHandler(r *http.Request) (*AdvertRequest, error) {
	var advert AdvertRequest
	err := json.NewDecoder(r.Body).Decode(&advert)
	if err != nil {
		return nil, fmt.Errorf("failed to decode advert request: %v", err)
	}
	return &advert, nil
}

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
