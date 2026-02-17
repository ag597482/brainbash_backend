package http

import (
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"

	"brainbash_backend/config"
	controller "brainbash_backend/internal/controller/http"
	"brainbash_backend/internal/middleware"
)

var (
	router *gin.Engine
	once   sync.Once
)

// initEngine initializes the gin engine with middlewares.
// Sets gin to release mode in production environments.
func initEngine(cfg *config.AppConfig, middlewares ...gin.HandlerFunc) {
	once.Do(func() {
		env := os.Getenv("ENVIRONMENT")
		if env == "prd" || env == "production" {
			gin.SetMode(gin.ReleaseMode)
		}

		appName := cfg.StaticConfig.App.Name
		if appName == "" {
			log.Fatal("APP_NAME cannot be empty!")
		}

		router = gin.New()
		middlewares = append(middlewares, gin.Logger(), gin.Recovery())
		router.Use(middlewares...)
	})
}

// Init initializes the http framework and registers all routes.
func Init(cfg *config.AppConfig) {
	initEngine(cfg)

	controllers := controller.NewControllers(cfg)

	// Public routes (no auth required)
	router.GET("/health", controllers.HealthController.Health)

	// Protected routes (JWT auth required)
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware(cfg.StaticConfig.Auth.JWTSecret))
	{
		// Register protected routes here, for example:
		// authorized.GET("/users", controllers.UserController.List)
	}
}

// Instance returns the initialized gin engine.
func Instance() *gin.Engine {
	if router == nil {
		log.Fatal("Router not initialized. Call Init() first.")
	}
	return router
}

// ResetForTesting resets the global state for testing purposes.
func ResetForTesting() {
	router = nil
	once = sync.Once{}
}
