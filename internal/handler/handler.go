package handler

import (
	"file-parser/internal/config"
	"file-parser/internal/middleware"
	"file-parser/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "file-parser/docs"
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
func (h *Handler) InitHandlers(config *config.LoggerConfig, apiKeyHash []byte) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Is used logger in middleware
	router.Use(middleware.InitLogger(config))
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	v1.GET("/health", h.health)
	v1.Use(middleware.APIKeyAuth(apiKeyHash))
	{
		v1.GET("/get-data", h.getDataByUnitGUID)
		v1.GET("/get-errors", h.getErrorsByFilename)
	}

	return router
}
