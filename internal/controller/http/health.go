package controller

import "github.com/gin-gonic/gin"

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) Health(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Application is up!!!"})
}