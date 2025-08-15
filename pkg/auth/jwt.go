package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User interface {
	GetID() string
	GetEmail() string
	GetRole() string
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret             string
	expiryHours        int
	refreshExpiryHours int
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWT(secret string, expiryHours, refreshExpiryHours int) *JWT {
	return &JWT{
		secret:             secret,
		expiryHours:        expiryHours,
		refreshExpiryHours: refreshExpiryHours,
	}
}

func (j *JWT) GenerateTokenPair(user User) (*TokenPair, error) {
	accessToken, err := j.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := j.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWT) GenerateAccessToken(user User) (string, error) {
	claims := &Claims{
		UserID: user.GetID(),
		Email:  user.GetEmail(),
		Role:   user.GetRole(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expiryHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.GetID(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWT) GenerateRefreshToken(user User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.refreshExpiryHours))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   user.GetID(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}

func (j *JWT) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", jwt.ErrInvalidKey
	}
	return claims.Subject, nil
}