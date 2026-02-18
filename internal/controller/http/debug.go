package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"brainbash_backend/config"
	"brainbash_backend/internal/model/response"
	"brainbash_backend/internal/service"
)

// DebugController exposes debug-only endpoints (e.g. user lookup by id, JWT from email).
type DebugController struct {
	cfg         *config.AppConfig
	userService *service.UserService
}

// NewDebugController creates a new DebugController with full config.
func NewDebugController(cfg *config.AppConfig, userService *service.UserService) *DebugController {
	return &DebugController{
		cfg:         cfg,
		userService: userService,
	}
}

// GetUser returns user details by user_id (path param).
// GET /debug/users/:user_id
func (dc *DebugController) GetUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	user, err := dc.userService.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Debug GetUser: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id or user not found"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, response.UserInfo{
		UserID:  user.UserID.Hex(),
		Email:   user.Email,
		Name:    user.Name,
		Picture: user.Picture,
	})
}

// GenerateJWTRequest is the body for POST /debug/jwt.
type GenerateJWTRequest struct {
	Email string `json:"email" binding:"required"`
}

// GenerateJWT returns a JWT for the user with the given email (debug only).
// POST /debug/jwt  body: { "email": "user@example.com" }
func (dc *DebugController) GenerateJWT(c *gin.Context) {
	var req GenerateJWTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	user, err := dc.userService.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		log.Printf("Debug GenerateJWT FindByEmail: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lookup user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserID.Hex(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(dc.cfg.StaticConfig.Auth.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"user_id":      user.UserID.Hex(),
		"email":        user.Email,
	})
}
