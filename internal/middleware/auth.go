package middleware

import (
	"errors"
	"file-parser/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func APIKeyAuth(apiKeyHash []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		if apiKey == "" {
			err := errors.New("API_KEY is required")
			c.Error(err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Error{
				Err: err.Error(),
			})
			return
		}

		if bcrypt.CompareHashAndPassword(apiKeyHash, []byte(apiKey)) != nil {
			err := errors.New("invalid API key")
			c.Error(err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Error{
				Err: err.Error(),
			})
			return
		}
		c.Next()
	}
}
