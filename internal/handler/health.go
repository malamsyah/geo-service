package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malamsyah/geo-service/internal/dto"
)

func Health(c *gin.Context) {
	healthResponse := dto.HealthResponse{
		Status: "ok",
	}

	c.JSON(http.StatusOK, healthResponse)
}
