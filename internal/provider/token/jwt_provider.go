package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayeed1999/share-a-ride/internal/config"
	"github.com/sayeed1999/share-a-ride/internal/domain/models"
)

type Claims struct {
	UserID   string          `json:"user_id"`
	UserType models.UserType `json:"user_type"`
	jwt.RegisteredClaims
}

type Provider interface {
	GenerateAccessToken(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type jwtProvider struct {
	config config.JWTConfig
}

func NewJWTProvider(config config.JWTConfig) Provider {
	return &jwtProvider{config: config}
}

func (p *jwtProvider) GenerateAccessToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   user.ID,
		UserType: user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(p.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.config.SecretKey))
}

func (p *jwtProvider) GenerateRefreshToken(user *models.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   user.ID,
		UserType: user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(p.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.config.SecretKey))
}

func (p *jwtProvider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
