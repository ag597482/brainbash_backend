package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"brainbash_backend/internal/model/entity"
	"brainbash_backend/internal/service"
)

// DashboardController handles public dashboard (leaderboard) API.
type DashboardController struct {
	dashboardService *service.DashboardService
}

// NewDashboardController creates a new DashboardController.
func NewDashboardController(dashboardService *service.DashboardService) *DashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
	}
}

// GetDashboard handles GET /api/dashboard. Returns top 10 scores per game type (public).
func (dc *DashboardController) GetDashboard(c *gin.Context) {
	d, err := dc.dashboardService.GetDashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load dashboard"})
		return
	}

	// Response: { processing_speed: [...], working_memory: [...], ... }
	out := map[string][]entity.DashboardEntry{
		"processing_speed":   d.ProcessingSpeed,
		"working_memory":     d.WorkingMemory,
		"logical_reasoning":  d.LogicalReasoning,
		"math_reasoning":     d.MathReasoning,
		"reflex_time":        d.ReflexTime,
		"attention_control":  d.AttentionControl,
	}
	// Ensure each key exists so client gets consistent shape
	if out["processing_speed"] == nil {
		out["processing_speed"] = []entity.DashboardEntry{}
	}
	if out["working_memory"] == nil {
		out["working_memory"] = []entity.DashboardEntry{}
	}
	if out["logical_reasoning"] == nil {
		out["logical_reasoning"] = []entity.DashboardEntry{}
	}
	if out["math_reasoning"] == nil {
		out["math_reasoning"] = []entity.DashboardEntry{}
	}
	if out["reflex_time"] == nil {
		out["reflex_time"] = []entity.DashboardEntry{}
	}
	if out["attention_control"] == nil {
		out["attention_control"] = []entity.DashboardEntry{}
	}

	c.JSON(http.StatusOK, out)
}
