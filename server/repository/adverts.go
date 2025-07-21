package repository

import "github.com/tousart/marketplace/models"

type AdvertsRepo interface {
	PostAdvert(advert *models.Advert) error
	GetAdvertsForward(req *models.FeedRequest) ([]models.FeedResponse, error)
	GetAdvertsBackward(req *models.FeedRequest) ([]models.FeedResponse, error)
}
