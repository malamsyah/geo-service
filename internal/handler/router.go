package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/malamsyah/geo-service/internal/middleware"
	"github.com/malamsyah/geo-service/internal/repository"
	"github.com/malamsyah/geo-service/internal/service"
	"github.com/malamsyah/geo-service/pkg/config"
	"gorm.io/gorm"
)

func SetupRouter(conf *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(middleware.JSONLoggerMiddleware())
	r.GET("/health", Health)

	// Setup geometry handler
	pointRepository := repository.NewPointRepository(db)
	contourRepository := repository.NewContourRepository(db)
	geometryService := service.NewGeometryService(pointRepository, contourRepository)
	geometryHandler := NewGeometryHandler(geometryService, conf.Host)

	defaultGroup := r.Group("/")
	geometryHandler.RegisterRoutes(defaultGroup)

	return r
}
