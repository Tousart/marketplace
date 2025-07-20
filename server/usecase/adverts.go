package usecase

type AdvertsService interface {
	PostAdvert(userID, title, text, url, price string) (string, string, error)
}
