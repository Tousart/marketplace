package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/marketplace/API/http/types"
	"github.com/tousart/marketplace/models"
	"github.com/tousart/marketplace/usecase"
)

type Adverts struct {
	advertsService usecase.AdvertsService
}

func NewAdvertsHandler(service usecase.AdvertsService) *Adverts {
	return &Adverts{advertsService: service}
}

func (a *Adverts) postAdvert(w http.ResponseWriter, r *http.Request) {
	advertReq, err := types.CreateAdvertRequestHandler(r)
	if err != nil {
		log.Printf("failed to create request advert handler: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(models.AuthKey).(string)
	if !ok {
		http.Error(w, "user id is empty", http.StatusBadRequest)
		return
	}

	advertID, date, err := a.advertsService.PostAdvert(userID,
		advertReq.Title,
		advertReq.Text,
		advertReq.URL,
		advertReq.Price)
	if err != nil {
		log.Printf("failed to post advert: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	advertResp := types.AdvertResponse{
		AdvertID: advertID,
		UserID:   userID,
		Title:    advertReq.Title,
		Date:     date,
	}
	if err := json.NewEncoder(w).Encode(&advertResp); err != nil {
		log.Printf("failed to encode response: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (a *Adverts) getAdvertsForward(w http.ResponseWriter, r *http.Request) {
	feedReq, err := types.CreateAdvertsFeedHandler(r)
	if err != nil {
		log.Printf("failed to create request adverts feed handler: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(models.AuthKey).(string)

	feedResp, err := a.advertsService.GetAdvertsForward(userID,
		feedReq.TypeSort,
		feedReq.FieldSort,
		feedReq.MinPrice,
		feedReq.MaxPrice,
		feedReq.LastVal,
		feedReq.LastID)
	if err != nil {
		log.Printf("failed to get adverts (forward): %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&feedResp); err != nil {
		log.Printf("failed to encode response (forward): %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (a *Adverts) getAdvertsBackward(w http.ResponseWriter, r *http.Request) {
	feedReq, err := types.CreateAdvertsFeedHandler(r)
	if err != nil {
		log.Printf("failed to create request adverts feed handler: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(models.AuthKey).(string)

	feedResp, err := a.advertsService.GetAdvertsBackward(userID,
		feedReq.TypeSort,
		feedReq.FieldSort,
		feedReq.MinPrice,
		feedReq.MaxPrice,
		feedReq.LastVal,
		feedReq.LastID)
	if err != nil {
		log.Printf("failed to get adverts (forward): %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&feedResp); err != nil {
		log.Printf("failed to encode response (forward): %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (a *Adverts) WithAdvertsHandlers(r chi.Router) {
	r.Route("/advert", func(r chi.Router) {
		r.Post("/", a.postAdvert)
		r.Get("/forward", a.getAdvertsForward)
		r.Get("/backward", a.getAdvertsBackward)
	})
}
