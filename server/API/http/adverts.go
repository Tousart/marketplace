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

	userID := r.Context().Value(models.AuthKey).(string)

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

func (a *Adverts) WithAdvertsHandlers(r chi.Router) {
	r.Post("/adverts", a.postAdvert)
}
