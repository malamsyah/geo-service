package service

import (
	"github.com/malamsyah/geo-service/internal/constants"
	"github.com/malamsyah/geo-service/internal/models"
	"github.com/malamsyah/geo-service/internal/repository"
)

type GeometryService interface {
	IsValidPoint(point *models.Point) bool
	CreatePoint(point *models.Point) error
	GetPoints(offset, limit int) ([]models.Point, error)
	GetPointByID(id uint) (*models.Point, error)
	IsValidContour(Contour *models.Contour) bool
	CreateContour(Contour *models.Contour) error
	GetContours(offset, limit int) ([]models.Contour, error)
	GetContourByID(id uint) (*models.Contour, error)
	UpdateContour(Contour *models.Contour) error
	DeleteContour(id uint) error

	// Advanced Query
	GetPointsByContourID(contourID uint) ([]models.Point, error)
	GetContoursIntersectArea(contourIDA, contourIDB uint) ([]models.Contour, error)
}

type GeometryServiceImpl struct {
	pointRepo   repository.PointRepository
	contourRepo repository.ContourRepository
}

func NewGeometryService(pointRepo repository.PointRepository, contourRepo repository.ContourRepository) GeometryService {
	return &GeometryServiceImpl{pointRepo, contourRepo}
}

func (s *GeometryServiceImpl) IsValidPoint(point *models.Point) bool {
	return point.Data.Validate() == nil
}

func (s *GeometryServiceImpl) CreatePoint(point *models.Point) error {
	if !s.IsValidPoint(point) {
		return constants.ErrInvalidPoint
	}

	return s.pointRepo.CreatePoint(point)
}

func (s *GeometryServiceImpl) GetPoints(offset, limit int) ([]models.Point, error) {
	return s.pointRepo.GetPoints(offset, limit)
}

func (s *GeometryServiceImpl) GetPointByID(id uint) (*models.Point, error) {
	return s.pointRepo.GetPointByID(id)
}

func (s *GeometryServiceImpl) IsValidContour(contour *models.Contour) bool {
	return contour.Data.Validate() == nil
}

func (s *GeometryServiceImpl) CreateContour(contour *models.Contour) error {
	if !s.IsValidContour(contour) {
		return constants.ErrInvalidContours
	}

	return s.contourRepo.CreateContour(contour)
}

func (s *GeometryServiceImpl) GetContours(offset, limit int) ([]models.Contour, error) {
	return s.contourRepo.GetContours(offset, limit)
}

func (s *GeometryServiceImpl) GetContourByID(id uint) (*models.Contour, error) {
	return s.contourRepo.GetContourByID(id)
}

func (s *GeometryServiceImpl) UpdateContour(contour *models.Contour) error {
	if !s.IsValidContour(contour) {
		return constants.ErrInvalidContours
	}

	return s.contourRepo.UpdateContour(contour)
}

func (s *GeometryServiceImpl) DeleteContour(id uint) error {
	return s.contourRepo.DeleteContour(id)
}

func (s *GeometryServiceImpl) GetPointsByContourID(contourID uint) ([]models.Point, error) {
	return s.pointRepo.GetPointsByContourID(contourID)
}

func (s *GeometryServiceImpl) GetContoursIntersectArea(contourIDA, contourIDB uint) ([]models.Contour, error) {
	_, err := s.contourRepo.GetContourByID(contourIDA)
	if err != nil {
		return nil, err
	}

	_, err = s.contourRepo.GetContourByID(contourIDB)
	if err != nil {
		return nil, err
	}

	resutls := make([]models.Contour, 0)
	geos, err := s.contourRepo.GetContoursIntersectArea(contourIDA, contourIDB)
	if err != nil {
		return nil, err
	}

	for _, c := range geos {
		if c.Data.IsPolygon() && len(c.Data.PolygonCoordinates) != 0 {
			resutls = append(resutls, c)
		}

		if c.Data.IsMultiPolygon() {
			for _, p := range c.Data.MultiPolygonCoordinates {
				resutls = append(resutls, models.Contour{Data: models.Geometry{Type: "Polygon", PolygonCoordinates: p}})
			}
		}
	}

	return resutls, nil
}
