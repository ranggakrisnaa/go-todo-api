package config

import (
	"errors"
	"time"

	"go-todo-api/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	JwtKey string
	JwtExp int
}

type Claims struct {
	UserID uint `json:"id"`
	jwt.RegisteredClaims
}

func NewJwtConfig(cfg *JwtConfig) (*JwtConfig, error) {
	return &JwtConfig{
		JwtKey: cfg.JwtKey,
		JwtExp: cfg.JwtExp,
	}, nil
}

func (c *JwtConfig) CreateToken(user *entity.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(c.JwtExp) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(c.JwtKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (c *JwtConfig) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(c.JwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
