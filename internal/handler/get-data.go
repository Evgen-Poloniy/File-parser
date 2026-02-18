package handler

import (
	"errors"
	"file-parser/internal/dto"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getDataByUnitGUID godoc
// @Summary Get parsed data by Unit GUID
// @Description Returns paginated parsed data for a given unit GUID
// @Tags data
// @Produce json
// @Param unit_guid query string true "Unit GUID"
// @Param limit query int false "Limit of records (1-1000)" default(1000) minimum(1) maximum(1000)
// @Param page query int false "Page number (>=1)" default(1) minimum(1)
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Error "Invalid request parameters"
// @Failure 500 {object} dto.Error "Internal server error"
// @Security ApiKeyAuth
// @Router /get-data [get]
func (h *Handler) getDataByUnitGUID(c *gin.Context) {
	ctx := c.Request.Context()

	unitGUID := c.Query("unit_guid")
	limitStr := c.DefaultQuery("limit", "1000")
	pageStr := c.DefaultQuery("page", "1")

	if unitGUID == "" {
		err := errors.New("value of query parameter unit_guid must not be empty")
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

	response, err := h.service.GetDataByUnitGUID(ctx, unitGUID, limit, page)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, dto.Error{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
