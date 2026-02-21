package controller

import (
	"brainbash_backend/config"
	appMongo "brainbash_backend/internal/mongo"
	"brainbash_backend/internal/repository"
	"brainbash_backend/internal/scoring"
	"brainbash_backend/internal/service"
	"strings"
)

// Controllers handles dependency injection in a centralized place.
type Controllers struct {
	HealthController    *HealthController
	AuthController      *AuthController
	DebugController     *DebugController
	ScoreController     *ScoreController
	DashboardController *DashboardController
	CleanupController   *CleanupController
}

func NewControllers(cfg *config.AppConfig) *Controllers {
	googleAuthService := service.NewGoogleAuthService(splitTrim(cfg.StaticConfig.Auth.GoogleClientID, ","))

	userRepo := repository.NewUserRepository(appMongo.GetDatabase())
	userService := service.NewUserService(userRepo)

	scorer := scoring.NewScorer()
	scoreRepo := repository.NewScoreRepository(appMongo.GetDatabase())
	dashboardRepo := repository.NewDashboardRepository(appMongo.GetDatabase())
	dashboardService := service.NewDashboardService(dashboardRepo, userService)
	scoreService := service.NewScoreService(scoreRepo, scorer, dashboardService)
	cleanupService := service.NewCleanupService(scoreRepo, dashboardRepo)

	return &Controllers{
		HealthController:    NewHealthController(),
		AuthController:      NewAuthController(googleAuthService, userService, cfg.StaticConfig.Auth.JWTSecret),
		DebugController:     NewDebugController(cfg, userService),
		ScoreController:     NewScoreController(scorer, scoreService),
		DashboardController: NewDashboardController(dashboardService),
		CleanupController:   NewCleanupController(cleanupService),
	}
}

// splitTrim splits s by sep and trims each element; returns nil for empty s.
func splitTrim(s, sep string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
