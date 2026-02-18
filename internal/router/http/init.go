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

	// CORS must be registered before any routes
	router.Use(middleware.CORSMiddleware())

	controllers := controller.NewControllers(cfg)

	// Public routes (no auth required)
	router.GET("/health", controllers.HealthController.Health)
	router.POST("/auth/google", controllers.AuthController.GoogleLogin)

	// Debug routes
	router.GET("/debug/users/:user_id", controllers.DebugController.GetUser)
	router.POST("/debug/jwt", controllers.DebugController.GenerateJWT)
	router.POST("/score", controllers.ScoreController.Calculate)

	// Protected routes (JWT auth required)
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware(cfg.StaticConfig.Auth.JWTSecret))
	{
		authorized.GET("/auth/me", controllers.AuthController.Me)
		authorized.POST("/api/game/result", controllers.ScoreController.GameResult)
		authorized.GET("/api/user/stats", controllers.ScoreController.UserStats)
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
