package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"brainbash_backend/internal/service"
)

const dateLayout = "02-01-2006" // dd-mm-yyyy

// CleanupController handles admin cleanup of scores and dashboard by date range.
type CleanupController struct {
	cleanupService *service.CleanupService
}

// NewCleanupController creates a new CleanupController.
func NewCleanupController(cleanupService *service.CleanupService) *CleanupController {
	return &CleanupController{
		cleanupService: cleanupService,
	}
}

// CleanupByDateRange handles DELETE /api/admin/cleanup?start_date=dd-mm-yyyy&end_date=dd-mm-yyyy.
// Removes sessions (from scores) and dashboard entries with timestamp in [start_date, end_date].
func (cc *CleanupController) CleanupByDateRange(c *gin.Context) {
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")
	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query params start_date and end_date are required (format: dd-mm-yyyy)"})
		return
	}

	start, err := time.Parse(dateLayout, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be dd-mm-yyyy"})
		return
	}
	end, err := time.Parse(dateLayout, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be dd-mm-yyyy"})
		return
	}

	// Use start of start day and end of end day (UTC)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, time.UTC)
	if end.Before(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be >= start_date"})
		return
	}

	scoresUpdated, err := cc.cleanupService.CleanupByDateRange(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cleanup failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "cleanup completed; sessions removed and avg_score/high_score/overall_score recomputed for affected users",
		"start_date":     startStr,
		"end_date":       endStr,
		"scores_updated": scoresUpdated,
	})
}
