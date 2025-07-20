package repository

import "github.com/tousart/marketplace/models"

type AdvertsRepo interface {
	PostAdvert(advert *models.Advert) error
}
