package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"brainbash_backend/internal/model/request"
	"brainbash_backend/internal/model/response"
	"brainbash_backend/internal/service"
)

type AuthController struct {
	googleAuthService *service.GoogleAuthService
	jwtSecret         string
}

func NewAuthController(googleAuthService *service.GoogleAuthService, jwtSecret string) *AuthController {
	return &AuthController{
		googleAuthService: googleAuthService,
		jwtSecret:         jwtSecret,
	}
}

// GoogleLogin handles POST /auth/google.
// Accepts either id_token (from mobile) or access_token (from web).
func (ac *AuthController) GoogleLogin(c *gin.Context) {
	var req request.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token or access_token is required"})
		return
	}

	var googleUser *service.GoogleUserInfo
	var err error

	switch {
	case req.IDToken != "":
		googleUser, err = ac.googleAuthService.VerifyIDToken(req.IDToken)
	case req.AccessToken != "":
		googleUser, err = ac.googleAuthService.VerifyAccessToken(req.AccessToken)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token or access_token is required"})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Google token: " + err.Error()})
		return
	}

	// Generate app JWT with user claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     googleUser.Sub,
		"email":   googleUser.Email,
		"name":    googleUser.Name,
		"picture": googleUser.Picture,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(ac.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		AccessToken: tokenString,
		User: response.UserInfo{
			ID:        googleUser.Sub,
			Email:     googleUser.Email,
			Name:      googleUser.Name,
			Picture:   googleUser.Picture,
			FirstName: googleUser.GivenName,
			LastName:  googleUser.FamilyName,
		},
	})
}

// Me handles GET /auth/me.
// Returns the current authenticated user's info from JWT claims.
func (ac *AuthController) Me(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No claims found"})
		return
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	c.JSON(http.StatusOK, response.UserInfo{
		ID:      getString(mapClaims, "sub"),
		Email:   getString(mapClaims, "email"),
		Name:    getString(mapClaims, "name"),
		Picture: getString(mapClaims, "picture"),
	})
}

func getString(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
