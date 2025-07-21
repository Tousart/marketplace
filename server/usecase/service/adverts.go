package service

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	if err := isValidTitle(title); err != nil {
		return "", "", fmt.Errorf("invalid title: %v", err)
	}

	if err := isValidText(text); err != nil {
		return "", "", fmt.Errorf("invalid text: %v", err)
	}

	if err := isValidURL(url); err != nil {
		return "", "", fmt.Errorf("invalid: url: %v", err)
	}

	if err := isValidPrice(price); err != nil {
		return "", "", fmt.Errorf("invalid price: %v", err)
	}

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
	advertString := strconv.Itoa(advertID)

	return advertString, dateString, nil
}

func (as *AdvertsService) GetAdvertsForward(userID, typeSort, fieldSort, minPriceStr, maxPriceStr, lastVal, lastIDStr string) ([]models.FeedResponse, error) {
	if err := isValidTypeSort(typeSort); err != nil {
		return nil, fmt.Errorf("invalid sort: %v", err)
	}

	if err := isValidFieldSort(fieldSort); err != nil {
		return nil, fmt.Errorf("invalid field sort: %v", err)
	}

	minPrice, maxPrice, err := isValidPriceRange(minPriceStr, maxPriceStr)
	if err != nil {
		return nil, fmt.Errorf("invalid price range: %v", err)
	}

	if lastVal == "" {
		if fieldSort == "date" {
			lastVal = "01-01-0001 00:00:00"
		} else {
			lastVal = "0"
		}
	}

	var lastID int
	if lastIDStr == "" {
		lastID = 0
	} else {
		lastID, err = strconv.Atoi(lastIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid last id: %v", err)
		}
	}

	var feed []models.FeedResponse
	feed, err = as.advertsRepo.GetAdvertsForward(&models.FeedRequest{
		TypeSort:  typeSort,
		FieldSort: fieldSort,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
		LastVal:   lastVal,
		LastID:    lastID,
		UserID:    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get forward error: %v", err)
	}

	return feed, err
}

func (as *AdvertsService) GetAdvertsBackward(userID, typeSort, fieldSort, minPriceStr, maxPriceStr, lastVal, lastIDStr string) ([]models.FeedResponse, error) {
	if err := isValidTypeSort(typeSort); err != nil {
		return nil, fmt.Errorf("invalid sort: %v", err)
	}

	if err := isValidFieldSort(fieldSort); err != nil {
		return nil, fmt.Errorf("invalid field sort: %v", err)
	}

	minPrice, maxPrice, err := isValidPriceRange(minPriceStr, maxPriceStr)
	if err != nil {
		return nil, fmt.Errorf("invalid price range: %v", err)
	}

	if lastVal == "" {
		if fieldSort == "date" {
			lastVal = "01-01-0001 00:00:00"
		} else {
			lastVal = "0"
		}
	}

	var lastID int
	if lastIDStr == "" {
		lastID = 0
	} else {
		lastID, err = strconv.Atoi(lastIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid last id: %v", err)
		}
	}

	var feed []models.FeedResponse
	feed, err = as.advertsRepo.GetAdvertsBackward(&models.FeedRequest{
		TypeSort:  typeSort,
		FieldSort: fieldSort,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
		LastVal:   lastVal,
		LastID:    lastID,
		UserID:    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("get forward error: %v", err)
	}

	return feed, err
}

func generateAdvertID(mu *sync.Mutex) (int, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile("/resources/id.txt")
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %v", err)
	}

	id, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("failed to convert count to int: %v", err)
	}

	newID := strconv.Itoa(id + 1)
	err = os.WriteFile("/resources/id.txt", []byte(newID), 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to write into file: %v", err)
	}

	return id, nil
}

func isValidTitle(title string) error {
	if len(title) < 3 {
		return errors.New("short title")
	} else if len(title) > 32 {
		return errors.New("long title")
	}

	regTrue := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	regFalse := regexp.MustCompile(`\-{2,}`)
	if !regTrue.MatchString(title) || regFalse.MatchString(title) {
		return errors.New("invalid title format")
	}

	return nil
}

func isValidText(text string) error {
	if len(text) < 3 {
		return errors.New("short text")
	} else if len(text) > 2000 {
		return errors.New("long text")
	}

	return nil
}

func isValidURL(url string) error {
	if !strings.HasPrefix(url, "htt") && !strings.HasPrefix(url, "ht") {
		return errors.New("invalid protocol")
	}

	if !strings.HasSuffix(url, "png") && !strings.HasSuffix(url, "jpg") {
		return errors.New("invalid format")
	}

	return nil
}

func isValidPrice(price string) error {
	if len(price) > 11 {
		return errors.New("so big price")
	}

	regTrue := regexp.MustCompile(`^[1-9][0-9]*$`)
	if !regTrue.MatchString(price) {
		return errors.New("invalid price")
	}
	return nil
}

func isValidTypeSort(typeSort string) error {
	if typeSort != "ascending" && typeSort != "descending" {
		return errors.New("type sort must be ascending or descending")
	}
	return nil
}

func isValidFieldSort(fieldSort string) error {
	if fieldSort != "date" && fieldSort != "price" {
		return errors.New("type sort must be date or price")
	}
	return nil
}

func isValidPriceRange(minPriceStr, maxPriceStr string) (int, int, error) {
	var minPrice int
	if minPriceStr == "" {
		minPrice = 0
	} else {
		var err error
		minPrice, err = strconv.Atoi(minPriceStr)
		if err != nil {
			return 0, 0, fmt.Errorf("incorrect min price format")
		}
	}

	var maxPrice int
	if maxPriceStr == "" {
		maxPrice = 1000000000
	} else {
		var err error
		maxPrice, err = strconv.Atoi(maxPriceStr)
		if err != nil {
			return 0, 0, fmt.Errorf("incorrect max price format")
		}
	}

	return minPrice, maxPrice, nil
}
