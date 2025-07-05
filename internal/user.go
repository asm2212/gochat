package internal

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	rdb *redis.Client
}

func NewUserService(rdb *redis.Client) *UserService {
	return &UserService{rdb: rdb}
}

const userPrefix = "user:"

var jwtSecret = []byte("supersecret") // In prod, use env var

// Register a new user
func (s *UserService) Register(ctx context.Context, username, password string) error {
	exists, _ := s.rdb.Exists(ctx, userPrefix+username).Result()
	if exists == 1 {
		return errors.New("username already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.rdb.HSet(ctx, userPrefix+username, "username", username, "password", string(hash)).Err()
}

// Login: checks password, returns JWT token if valid
func (s *UserService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.rdb.HGetAll(ctx, userPrefix+username).Result()
	if err != nil || len(user) == 0 {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user["password"]), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	// JWT: username, exp 24h
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// ParseToken extracts username from JWT
func (s *UserService) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("invalid token username")
	}
	return username, nil
}
