package handler

import (
	"errors"
	"file-parser/internal/dto"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getErrorsByFilename godoc
// @Summary Get parsing errors by filename
// @Description Returns paginated errors for a given filename
// @Tags errors
// @Produce json
// @Param filename query string true "Filename"
// @Param limit query int false "Limit of records (1-1000)" default(1000) minimum(1) maximum(1000)
// @Param page query int false "Page number (>=1)" default(1) minimum(1)
// @Success 200 {object} dto.FileErrorInfo
// @Failure 400 {object} dto.Error "Invalid request parameters"
// @Failure 500 {object} dto.Error "Internal server error"
// @Security ApiKeyAuth
// @Router /get-errors [get]
func (h *Handler) getErrorsByFilename(c *gin.Context) {
	ctx := c.Request.Context()

	filename := c.Query("filename")
	limitStr := c.DefaultQuery("limit", "1000")
	pageStr := c.DefaultQuery("page", "1")

	if filename == "" {
		err := errors.New("value of query parameter filename must not be empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, dto.Error{Err: err.Error()})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.Error(fmt.Errorf("bad conversion of query parameter limit: %w", err))
		c.JSON(http.StatusBadRequest, dto.Error{Err: err.Error()})
		return
	}

	if limit < 1 {
		err := errors.New("value of query parameter limit must be grater then 1")
		c.Error(err)
		c.JSON(http.StatusBadRequest, dto.Error{Err: err.Error()})
		return
	}

	if limit > 1000 {
		err := errors.New("value of query parameter limit must be less then 1000")
		c.Error(err)
		c.JSON(http.StatusBadRequest, dto.Error{Err: err.Error()})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Error(fmt.Errorf("bad conversion of query parameter pageStr: %w", err))
		c.JSON(http.StatusBadRequest, dto.Error{Err: err.Error()})
		return
	}

	if page < 1 {
		err := errors.New("value of query parameter page must be grater then 1")
		c.Error(err)
		c.JSON(http.StatusBadRequest, dto.Error{Err: err.Error()})
		return
	}

	response, err := h.service.GetFileErrorsByFilename(ctx, filename, limit, page)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.Error{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
