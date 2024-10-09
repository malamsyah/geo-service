package repository

import (
	"fmt"

	"github.com/malamsyah/geo-service/internal/models"
	"gorm.io/gorm"
)

type PointRepository interface {
	CreatePoint(point *models.Point) error
	GetPointByID(id uint) (*models.Point, error)
	GetPointsByContourID(contourID uint) ([]models.Point, error)
	GetPoints(offset, limit int) ([]models.Point, error)
	UpdatePoint(point *models.Point) error
	DeletePoint(id uint) error
}

type PointRepositoryImpl struct {
	db *gorm.DB
}

type filter struct {
	Offset int
	Limit  int
	ID     uint
}

func NewPointRepository(db *gorm.DB) PointRepository {
	return &PointRepositoryImpl{db}
}

func (r *PointRepositoryImpl) CreatePoint(point *models.Point) error {
	return r.db.Create(point).Error
}

func (r *PointRepositoryImpl) GetPointByID(id uint) (*models.Point, error) {
	point := new(models.Point)
	query, params := r.getPointQuery(filter{ID: id})

	err := r.db.Raw(query, params...).Scan(&point).Error
	if err != nil {
		return nil, err
	}

	if point.ID == uint(0) {
		return nil, gorm.ErrRecordNotFound
	}

	return point, nil
}

func (r *PointRepositoryImpl) GetPoints(offset, limit int) ([]models.Point, error) {
	var points []models.Point
	query, params := r.getPointQuery(filter{Offset: offset, Limit: limit})

	err := r.db.Raw(query, params...).Scan(&points).Error
	if err != nil {
		return nil, err
	}

	return points, nil
}

func (r *PointRepositoryImpl) UpdatePoint(point *models.Point) error {
	return r.db.Save(point).Error
}

func (r *PointRepositoryImpl) DeletePoint(id uint) error {
	return r.db.Delete(&models.Point{}, id).Error
}

func (r *PointRepositoryImpl) getPointQuery(f filter) (string, []any) {
	params := make([]any, 0)
	query := "SELECT id, ST_AsGeoJSON(data) AS data FROM points"
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

func (r *PointRepositoryImpl) GetPointsByContourID(contourID uint) ([]models.Point, error) {
	points := make([]models.Point, 0)
	query := "SELECT p.id, ST_AsGeoJSON(p.data) AS data FROM points p JOIN contours c ON ST_Within(p.data, c.data) WHERE c.id = ?"
	err := r.db.Raw(query, contourID).Scan(&points).Error
	if err != nil {
		return nil, err
	}

	return points, nil
}
