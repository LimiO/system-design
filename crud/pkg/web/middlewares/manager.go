package middlewares

import "onlinestore/internal/jwt"

type MiddlewareManager struct {
	tokenManager *jwt.TokenManager
}

func NewMiddlewareManager(token string) *MiddlewareManager {
	tokenManager := jwt.NewTokenManager(token)
	return &MiddlewareManager{
		tokenManager: tokenManager,
	}
}
