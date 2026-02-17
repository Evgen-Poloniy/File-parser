package handler

import (
	"file-parser/internal/config"
	"file-parser/internal/middleware"
	"file-parser/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler
type Handler struct {
	service service.Reader
}

// NewHandler creates a new Handler with the given service
func NewHandler(service service.Reader) *Handler {
	return &Handler{
		service: service,
	}
}

// InitHandlers initializes the HTTP handlers, routes, and middleware.
func (h *Handler) InitHandlers(config *config.LoggerConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Is used logger in middleware
	router.Use(middleware.InitLogger(config))
	router.Use(gin.Recovery())

	router.GET("/health", h.health)
	router.GET("/get-data", h.getDataByUnitGUID)
	router.GET("/get-errors", h.getErrorsByFilename)

	return router
}
