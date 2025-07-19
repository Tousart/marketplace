package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tousart/marketplace/API/http/types"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/marketplace/service"
)

type Users struct {
	authService service.AuthService
}

func NewUsersHandler(auth service.AuthService) *Users {
	return &Users{authService: auth}
}

func (u *Users) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	userReq, err := types.CreateRequestHandler(r)
	if err != nil {
		log.Printf("bad request: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	token, err := u.authService.Login(userReq.Login, userReq.Password)
	if err != nil {
		log.Printf("internal error: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&types.LoginResponse{Token: token}); err != nil {
		log.Printf("failed to encode response: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (u *Users) postRegisterHandler(w http.ResponseWriter, r *http.Request) {
	userReq, err := types.CreateRequestHandler(r)
	if err != nil {
		log.Printf("bad request: %v\n", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	userID, login, hashPassword, err := u.authService.Register(userReq.Login, userReq.Password)
	if err != nil {
		log.Printf("internal error: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	userResp := types.RegisterResponse{
		UserID:       userID,
		Login:        login,
		HashPassword: hashPassword,
	}

	if err := json.NewEncoder(w).Encode(&userResp); err != nil {
		log.Printf("failed to encode response: %v\n", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (u *Users) WithUsersHandlers(r chi.Router) {
	r.Post("/login", u.postLoginHandler)
	r.Post("/register", u.postRegisterHandler)
}
