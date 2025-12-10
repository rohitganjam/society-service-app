package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rohit/society-service-app/backend/internal/database"
	"github.com/rohit/society-service-app/backend/internal/utils"
)

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthStatus struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Time     string            `json:"time"`
	Services map[string]string `json:"services"`
}

func (h *HealthHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)

	// Check database
	if h.db != nil {
		if err := h.db.Health(ctx); err != nil {
			services["database"] = "unhealthy"
		} else {
			services["database"] = "healthy"
		}
	} else {
		services["database"] = "not_configured"
	}

	status := "healthy"
	for _, v := range services {
		if v == "unhealthy" {
			status = "unhealthy"
			break
		}
	}

	utils.RespondSuccess(c, http.StatusOK, HealthStatus{
		Status:   status,
		Version:  "1.0.0",
		Time:     time.Now().UTC().Format(time.RFC3339),
		Services: services,
	}, "")
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if h.db != nil {
		if err := h.db.Health(ctx); err != nil {
			utils.RespondError(c, http.StatusServiceUnavailable, "NOT_READY", "Database not ready", nil)
			return
		}
	}

	utils.RespondSuccess(c, http.StatusOK, gin.H{"ready": true}, "")
}
