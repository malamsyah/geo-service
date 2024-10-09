package repository

import (
	"fmt"

	"github.com/malamsyah/geo-service/internal/constants"
	"github.com/malamsyah/geo-service/internal/models"
	"gorm.io/gorm"
)

type ContourRepository interface {
	CreateContour(Contour *models.Contour) error
	GetContourByID(id uint) (*models.Contour, error)
	GetContours(offset, limit int) ([]models.Contour, error)
	UpdateContour(contour *models.Contour) error
	DeleteContour(id uint) error
	GetContoursIntersectArea(idA, idB uint) ([]models.Contour, error)
}

type ContourRepositoryImpl struct {
	db *gorm.DB
}

func NewContourRepository(db *gorm.DB) ContourRepository {
	return &ContourRepositoryImpl{db}
}

func (r *ContourRepositoryImpl) CreateContour(contour *models.Contour) error {
	return r.db.Create(contour).Error
}

func (r *ContourRepositoryImpl) GetContourByID(id uint) (*models.Contour, error) {
	contour := new(models.Contour)
	query, params := r.getContourQuery(filter{ID: id})

	err := r.db.Raw(query, params...).Scan(&contour).Error
	if err != nil {
		return nil, err
	}

	if contour.ID == uint(0) {
		return nil, constants.ErrContourNotFound
	}

	return contour, nil
}

func (r *ContourRepositoryImpl) GetContours(offset, limit int) ([]models.Contour, error) {
	var contours []models.Contour
	query, params := r.getContourQuery(filter{Offset: offset, Limit: limit})

	err := r.db.Raw(query, params...).Scan(&contours).Error
	if err != nil {
		return nil, err
	}

	return contours, nil
}

func (r *ContourRepositoryImpl) UpdateContour(contour *models.Contour) error {
	return r.db.Save(contour).Error
}

func (r *ContourRepositoryImpl) DeleteContour(id uint) error {
	return r.db.Delete(&models.Contour{}, id).Error
}

func (r *ContourRepositoryImpl) getContourQuery(f filter) (string, []any) {
	params := make([]any, 0)
	query := "SELECT id, ST_AsGeoJSON(data) AS data FROM contours"
	keyword := "WHERE"
	if f.ID != 0 {
		query += fmt.Sprintf(" %s id = ?", keyword)
		params = append(params, f.ID)
	} else {
		query += " ORDER BY id DESC OFFSET ? LIMIT ?"
		params = append(params, f.Offset, f.Limit)
	}

	return query, params
}

func (r *ContourRepositoryImpl) GetContoursIntersectArea(idA, idB uint) ([]models.Contour, error) {
	contours := make([]models.Contour, 0)
	query := "SELECT ST_AsGeoJSON(ST_Intersection(ca.data, cb.data)) AS data FROM contours ca, contours cb WHERE ca.id = ? AND cb.id = ?"
	err := r.db.Raw(query, idA, idB).Scan(&contours).Error
	if err != nil {
		return nil, err
	}

	return contours, nil
}
