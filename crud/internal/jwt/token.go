package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{
		secret: secret,
	}
}

func (t *TokenManager) CreateToken(username string) (string, error) {
	maker := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": username,
		})
	s, err := maker.SignedString([]byte(t.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign string: %v", err)
	}
	return s, nil
}

func (t *TokenManager) GetClaims(token string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected data signing method: %v", token.Header["alg"])
		}
		return []byte(t.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse with claims: %v", err)
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	return claims, nil
}
