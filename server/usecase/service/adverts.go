package service

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/tousart/marketplace/models"
	"github.com/tousart/marketplace/repository"
)

type AdvertsService struct {
	advertsRepo repository.AdvertsRepo
	mutex       *sync.Mutex
}

func NewAdvertsService(repo repository.AdvertsRepo, mu *sync.Mutex) *AdvertsService {
	return &AdvertsService{
		advertsRepo: repo,
		mutex:       mu,
	}
}

func (as *AdvertsService) PostAdvert(userID, title, text, url, price string) (string, string, error) {
	advertID, err := generateAdvertID(as.mutex)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate advert id: %v", err)
	}

	date := time.Now()

	intPrice, err := strconv.Atoi(price)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert price to int: %v", err)
	}

	advert := models.Advert{
		AdvertID: advertID,
		UserID:   userID,
		Title:    title,
		Text:     text,
		URL:      url,
		Price:    intPrice,
		Date:     date,
	}

	if err := as.advertsRepo.PostAdvert(&advert); err != nil {
		return "", "", fmt.Errorf("failed to post advert into db: %v", err)
	}

	dateString := date.Format("02-01-2006 15:04:05")

	return advertID, dateString, nil
}

func generateAdvertID(mu *sync.Mutex) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile("./resources/id.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	number, err := strconv.Atoi(string(data))
	if err != nil {
		return "", fmt.Errorf("failed to convert count to int: %v", err)
	}

	id := strconv.Itoa(number + 1)
	err = os.WriteFile("./resources/id.txt", []byte(id), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write into file: %v", err)
	}

	return id, nil
}
