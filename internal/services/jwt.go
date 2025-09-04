package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret         []byte
	expiresMinutes int
}

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, expiresMinutes int) *JWTService {
	return &JWTService{secret: []byte(secret), expiresMinutes: expiresMinutes}
}

func (j *JWTService) Generate(userID uuid.UUID) (string, error) {
	expiresAt := time.Now().Add(time.Duration(j.expiresMinutes) * time.Minute)
	claims := &Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

func (j *JWTService) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
