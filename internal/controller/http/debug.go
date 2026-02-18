package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"brainbash_backend/internal/model/response"
	"brainbash_backend/internal/service"
)

// DebugController exposes debug-only endpoints (e.g. user lookup by id).
type DebugController struct {
	userService *service.UserService
}

// NewDebugController creates a new DebugController.
func NewDebugController(userService *service.UserService) *DebugController {
	return &DebugController{
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
