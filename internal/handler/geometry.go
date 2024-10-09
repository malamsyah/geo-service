package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/malamsyah/geo-service/internal/constants"
	"github.com/malamsyah/geo-service/internal/dto"
	"github.com/malamsyah/geo-service/internal/service"
	"github.com/malamsyah/geo-service/pkg/logger"
)

const (
	DefaultLimit  = 10
	DefaultOffset = 0
)

type GeometryHandler struct {
	geometryService service.GeometryService
	host            string
}

func NewGeometryHandler(geometryService service.GeometryService, host string) *GeometryHandler {
	return &GeometryHandler{geometryService, host}
}

func (h *GeometryHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/points", h.CreatePoint)
	r.GET("/points", h.GetPoints)
	r.POST("/contours", h.CreateContour)
	r.GET("/contours", h.GetContours)
	r.GET("/contours/:id", h.GetContourByID)
	r.PUT("/contours/:id", h.UpdateContour)
	r.DELETE("/contours/:id", h.DeleteContour)
	r.GET("/intersections", h.Intersect)
}

func (h *GeometryHandler) CreatePoint(c *gin.Context) {
	var req dto.CreatePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	point := req.ToModel()

	err := h.geometryService.CreatePoint(&point)
	if err != nil {
		logger.Errorf("Failed to create point: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, point)
}

func (h *GeometryHandler) GetPoints(c *gin.Context) {
	page, offset, limit, err := h.parseOffsetLimit(c)
	if err != nil {
		logger.Errorf("Failed to parse page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conourIDStr := c.Query("contour")
	if conourIDStr != "" {
		contourID, err := strconv.Atoi(conourIDStr)
		if err != nil {
			logger.Errorf("Failed to parse contour id: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		points, err := h.geometryService.GetPointsByContourID(uint(contourID))
		if err != nil {
			logger.Errorf("Failed to get points: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := dto.Response{
			Count:    len(points),
			Next:     h.buildNextURL("/points", page),
			Previous: h.buildPreviousURL("/points", page),
			Results:  points,
		}

		c.JSON(http.StatusOK, resp)
		return
	}

	points, err := h.geometryService.GetPoints(offset, limit)
	if err != nil {
		logger.Errorf("Failed to get points: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.Response{
		Count:    len(points),
		Next:     h.buildNextURL("/points", page),
		Previous: h.buildPreviousURL("/points", page),
		Results:  points,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *GeometryHandler) CreateContour(c *gin.Context) {
	var req dto.CreateContourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Contour := req.ToModel()

	err := h.geometryService.CreateContour(&Contour)
	if err != nil {
		logger.Errorf("Failed to create Contour: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, Contour)
}

func (h *GeometryHandler) GetContours(c *gin.Context) {
	page, offset, limit, err := h.parseOffsetLimit(c)
	if err != nil {
		logger.Errorf("Failed to parse page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contours, err := h.geometryService.GetContours(offset, limit)
	if err != nil {
		logger.Errorf("Failed to get contours: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dto.Response{
		Count:    len(contours),
		Next:     h.buildNextURL("/contours", page),
		Previous: h.buildPreviousURL("/contours", page),
		Results:  contours,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *GeometryHandler) GetContourByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("Failed to parse id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contour, err := h.geometryService.GetContourByID(uint(id))
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logger.Errorf("Failed to get contour: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contour)
}

func (h *GeometryHandler) UpdateContour(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("Failed to parse id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req dto.CreateContourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contour := req.ToModel()
	contour.ID = uint(id)

	err = h.geometryService.UpdateContour(&contour)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logger.Errorf("Failed to update contour: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contour)
}

func (h *GeometryHandler) DeleteContour(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("Failed to parse id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.geometryService.DeleteContour(uint(id))
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		logger.Errorf("Failed to delete contour: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *GeometryHandler) Intersect(c *gin.Context) {
	contourIDA, err := strconv.Atoi(c.Query("contour_1"))
	if err != nil {
		logger.Errorf("Failed to parse contourA: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contourIDB, err := strconv.Atoi(c.Query("contour_2"))
	if err != nil {
		logger.Errorf("Failed to parse contourB: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contours, err := h.geometryService.GetContoursIntersectArea(uint(contourIDA), uint(contourIDB))
	if err != nil {
		logger.Errorf("Failed to get contours: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contours)
}

func (h *GeometryHandler) parseOffsetLimit(c *gin.Context) (int, int, int, error) {
	pageStr := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, 0, err
	}

	offset := page * DefaultLimit
	limit := DefaultLimit

	return page, offset, limit, nil
}

func (h *GeometryHandler) buildNextURL(path string, page int) *string {
	res := h.host + path + "?page=" + fmt.Sprint(page+1)
	return &res
}

func (h *GeometryHandler) buildPreviousURL(path string, page int) *string {
	nextPage := page - 1
	if nextPage < 0 {
		return nil
	}

	res := h.host + path + "?page=" + fmt.Sprint(nextPage)
	return &res
}
