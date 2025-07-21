package usecase

import "github.com/tousart/marketplace/models"

type AdvertsService interface {
	PostAdvert(userID, title, text, url, price string) (string, string, error)
	GetAdvertsForward(userID, typeSort, fieldSort, minPriceStr, maxPriceStr, lastVal, lastIDStr string) ([]models.FeedResponse, error)
	GetAdvertsBackward(userID, typeSort, fieldSort, minPriceStr, maxPriceStr, lastVal, lastIDStr string) ([]models.FeedResponse, error)
}
