package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// health godoc
// @Summary Health check
// @Description Check service availability
// @Tags system
// @Success 200 {string} string "OK"
// @Router /health [get]
func (h *Handler) health(c *gin.Context) {
	c.Status(http.StatusOK)
}
