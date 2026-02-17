package controller

import "brainbash_backend/config"

// handles dependency injection in a centralized place
type Controllers struct {
	HealthController *HealthController
}

func NewControllers(cfg *config.AppConfig) *Controllers {
	return &Controllers{
		HealthController: NewHealthController(),
	}
}
