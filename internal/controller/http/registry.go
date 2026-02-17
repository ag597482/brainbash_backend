package controller

import (
	"brainbash_backend/config"
	"brainbash_backend/internal/service"
)

// handles dependency injection in a centralized place
type Controllers struct {
	HealthController *HealthController
	AuthController   *AuthController
}

func NewControllers(cfg *config.AppConfig) *Controllers {
	googleAuthService := service.NewGoogleAuthService(cfg.StaticConfig.Auth.GoogleClientID)

	return &Controllers{
		HealthController: NewHealthController(),
		AuthController:   NewAuthController(googleAuthService, cfg.StaticConfig.Auth.JWTSecret),
	}
}
