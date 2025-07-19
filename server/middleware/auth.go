package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/tousart/marketplace/service/auth"
)

type authKey string

const key authKey = "id"

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		if tokenString == header {
			http.Error(w, "invalid token format", http.StatusBadRequest)
			return
		}

		userID, err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Printf("invalid token: %v\n", err)
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), key, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
