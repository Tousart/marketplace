package auth

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tousart/marketplace/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	TOKEN_EXP  = int64(12 * time.Hour)
	SECRET_KEY = "ne-secret"
)

type AuthService struct {
	usersRepo repository.UsersRepo
}

func NewAuthService(repo repository.UsersRepo) *AuthService {
	return &AuthService{
		usersRepo: repo,
	}
}

func (as *AuthService) Login(login, password string) (string, error) {
	userID, err := as.usersRepo.Login(login, password)
	if err != nil {
		return "", fmt.Errorf("failed to login user: %v", err)
	}

	token, err := generateToken(userID, login)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func (as *AuthService) Register(login, password string) (string, string, string, error) {
	if err := isValidRegister(login, password); err != nil {
		return "", "", "", err
	}

	userID := generateID()

	hashPassword, err := hash(password)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to hash password: %v", err)
	}

	return userID, login, hashPassword, as.usersRepo.Register(userID, login, hashPassword)
}

func isValidRegister(login, password string) error {
	if len(login) < 4 {
		return errors.New("short login")
	} else if len(login) > 32 {
		return errors.New("long login")
	} else if len(password) < 8 {
		return errors.New("short password")
	} else if len(password) > 32 {
		return errors.New("long password")
	}

	regTrue := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	regFalse := regexp.MustCompile(`[._-]{2,}`)
	if !regTrue.MatchString(login) || regFalse.MatchString(login) {
		return errors.New("invalid login format")
	} else if !regTrue.MatchString(password) || regFalse.MatchString(password) {
		return errors.New("invalid password format")
	}

	return nil
}

func generateID() string {
	id := uuid.NewString()
	return id
}

func hash(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash error: %v", err)
	}
	return string(hashPassword), nil
}

func generateToken(userID, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    userID,
			"login": login,
			"exp":   TOKEN_EXP,
		})

	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid method")
		}
		return SECRET_KEY, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims format")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("user id not found")
	}

	return userID, nil
}
