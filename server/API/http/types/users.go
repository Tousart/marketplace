package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateUserRequestHandler(r *http.Request) (*UserRequest, error) {
	var user UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user request: %v", err)
	}
	return &user, nil
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	UserID       string `json:"user_id"`
	Login        string `json:"login"`
	HashPassword string `json:"hash_password"`
}
