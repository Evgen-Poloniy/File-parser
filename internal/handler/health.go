package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) health(c *gin.Context) {
	c.Status(http.StatusOK)
}
