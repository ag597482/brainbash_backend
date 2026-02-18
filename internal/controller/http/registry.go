package controller

import (
	"brainbash_backend/config"
	appMongo "brainbash_backend/internal/mongo"
	"brainbash_backend/internal/repository"
	"brainbash_backend/internal/scoring"
	"brainbash_backend/internal/service"
)

// Controllers handles dependency injection in a centralized place.
type Controllers struct {
	HealthController  *HealthController
	AuthController    *AuthController
	DebugController   *DebugController
	ScoringController *ScoringController
}

func NewControllers(cfg *config.AppConfig) *Controllers {
	googleAuthService := service.NewGoogleAuthService(cfg.StaticConfig.Auth.GoogleClientID)

	userRepo := repository.NewUserRepository(appMongo.GetDatabase())
	userService := service.NewUserService(userRepo)

	scorer := scoring.NewScorer()

	return &Controllers{
		HealthController:  NewHealthController(),
		AuthController:    NewAuthController(googleAuthService, userService, cfg.StaticConfig.Auth.JWTSecret),
		DebugController:   NewDebugController(cfg, userService),
		ScoringController: NewScoringController(scorer),
	}
}
