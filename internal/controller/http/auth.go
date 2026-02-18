package controller

import (
	"log"
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
	userService       *service.UserService
	jwtSecret         string
}

func NewAuthController(googleAuthService *service.GoogleAuthService, userService *service.UserService, jwtSecret string) *AuthController {
	return &AuthController{
		googleAuthService: googleAuthService,
		userService:       userService,
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

	// If user already exists by email, use stored details and do not create a new user
	persistedUser, err := ac.userService.FindByEmail(c.Request.Context(), googleUser.Email)
	if err != nil {
		log.Printf("Failed to find user by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lookup user"})
		return
	}
	if persistedUser == nil {
		// New user: create in MongoDB
		persistedUser, err = ac.userService.UpsertFromGoogleLogin(c.Request.Context(), googleUser)
		if err != nil {
			log.Printf("Failed to save user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
			return
		}
	}

	// Generate app JWT with user_id as the subject
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": persistedUser.UserID.Hex(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(ac.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Use stored user details for response (existing or newly created)
	c.JSON(http.StatusOK, response.LoginResponse{
		AccessToken: tokenString,
		User: response.UserInfo{
			UserID:    persistedUser.UserID.Hex(),
			Email:     persistedUser.Email,
			Name:      persistedUser.Name,
			Picture:   persistedUser.Picture,
			FirstName: googleUser.GivenName,
			LastName:  googleUser.FamilyName,
		},
	})
}

// Me handles GET /auth/me.
// Reads user_id (sub) from JWT claims, fetches user from DB.
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

	userID := getString(mapClaims, "sub")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing sub"})
		return
	}

	user, err := ac.userService.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
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

func getString(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
