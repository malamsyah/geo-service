package repository

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/malamsyah/geo-service/internal/db"
	"github.com/malamsyah/geo-service/internal/models"
	"github.com/malamsyah/geo-service/pkg/config"
)

type PointRepoTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestPointRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PointRepoTestSuite))
}

func setupTestDB(t *testing.T) *gorm.DB {
	conf := config.Load(path.Join("../", "../", ".env"))
	dbCon, err := db.ConnectPostgres(conf)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Migrate(dbCon)
	if err != nil {
		t.Fatal(err)
	}

	return dbCon
}

func (p *PointRepoTestSuite) SetupSuite() {
	p.db = setupTestDB(p.Suite.T())
}

func (p *PointRepoTestSuite) TestPointRepository_CreatePoint() {
	tx := p.db.Begin()
	repo := NewPointRepository(tx)
	p.Suite.T().Run("CreatePoint", func(t *testing.T) {
		tests := []struct {
			name    string
			point   *models.Point
			wantErr bool
		}{
			{
				name: "ValidPoint",
				point: &models.Point{
					Data: models.Geometry{
						Type:             "Point",
						PointCoordinates: [2]float64{125.6, 10.1},
					},
				},
				wantErr: false,
			},
			{
				name: "EmptyData",
				point: &models.Point{
					Data: models.Geometry{
						Type:             "Point",
						PointCoordinates: [2]float64{},
					},
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.CreatePoint(tt.point)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, tt.point.ID)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *PointRepoTestSuite) TestPointRepository_GetPointByID() {
	tx := p.db.Begin()
	repo := NewPointRepository(tx)

	examplePoint := &models.Point{
		Data: models.Geometry{
			Type:             "Point",
			PointCoordinates: [2]float64{125.6, 10.1},
		},
	}
	err := repo.CreatePoint(examplePoint)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	p.Suite.T().Run("GetPointByID", func(t *testing.T) {
		tests := []struct {
			name    string
			ID      uint
			wantErr bool
		}{
			{
				name:    "PointExists",
				ID:      examplePoint.ID,
				wantErr: false,
			},
			{
				name:    "PointNotExists",
				ID:      999999,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actualPoint, err := repo.GetPointByID(tt.ID)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, actualPoint.ID)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *PointRepoTestSuite) TestPointRepository_GetPoints() {
	tx := p.db.Begin()
	repo := NewPointRepository(tx)

	examplePoint := &models.Point{
		Data: models.Geometry{
			Type:             "Point",
			PointCoordinates: [2]float64{125.6, 10.1},
		},
	}
	err := repo.CreatePoint(examplePoint)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	p.Suite.T().Run("GetPoints", func(t *testing.T) {
		tests := []struct {
			name           string
			limit          int
			offset         int
			expectedLength int
			wantErr        bool
		}{
			{
				name:           "OnePoint",
				limit:          1,
				offset:         0,
				expectedLength: 1,
				wantErr:        false,
			},
			{
				name:           "ZeroLimit",
				limit:          0,
				offset:         0,
				expectedLength: 0,
				wantErr:        false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				points, err := repo.GetPoints(tt.offset, tt.limit)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedLength, len(points))
				}
			})
		}
	})

	tx.Rollback()
}

func (p *PointRepoTestSuite) TestPointRepository_UpdatePoint() {
	tx := p.db.Begin()
	repo := NewPointRepository(tx)

	examplePoint := &models.Point{
		Data: models.Geometry{
			Type:             "Point",
			PointCoordinates: [2]float64{125.6, 10.1},
		},
	}
	err := repo.CreatePoint(examplePoint)
	if err != nil {
		p.Suite.T().Fatal(err)
	}
	p.Suite.T().Run("UpdatePoint", func(t *testing.T) {
		tests := []struct {
			name    string
			point   *models.Point
			wantErr bool
		}{
			{
				name: "ValidPoint",
				point: &models.Point{
					ID: examplePoint.ID,
					Data: models.Geometry{
						Type:             "Point",
						PointCoordinates: [2]float64{125.6, 10.1},
					},
				},
				wantErr: false,
			},
			{
				name: "EmptyData",
				point: &models.Point{
					ID: examplePoint.ID,
					Data: models.Geometry{
						Type: "Point",
					},
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.UpdatePoint(tt.point)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, tt.point.ID)
					assert.Equal(t, examplePoint.ID, tt.point.ID)

					actualPoint, err := repo.GetPointByID(tt.point.ID)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tt.point.Data, actualPoint.Data)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *PointRepoTestSuite) TestPointRepository_DeletePoint() {
	tx := p.db.Begin()
	repo := NewPointRepository(tx)

	examplePoint := &models.Point{
		Data: models.Geometry{
			Type:             "Point",
			PointCoordinates: [2]float64{125.6, 10.1},
		},
	}
	err := repo.CreatePoint(examplePoint)
	if err != nil {
		p.Suite.T().Fatal(err)
	}
	p.Suite.T().Run("DeletePoint", func(t *testing.T) {
		tests := []struct {
			name    string
			ID      uint
			wantErr bool
		}{
			{
				name:    "PointExists",
				ID:      examplePoint.ID,
				wantErr: false,
			},
			{
				name:    "PointNotExists",
				ID:      999999,
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.DeletePoint(tt.ID)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					actualPoint, err := repo.GetPointByID(tt.ID)
					assert.Error(t, err)
					assert.Nil(t, actualPoint)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *PointRepoTestSuite) TestPointRepository_GetPointsByContourID() {
	tx := p.db.Begin()
	repo := NewPointRepository(tx)

	examplePoint := &models.Point{
		Data: models.Geometry{
			Type:             "Point",
			PointCoordinates: [2]float64{5.0, 5.0},
		},
	}
	err := repo.CreatePoint(examplePoint)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	contourRepo := NewContourRepository(tx)
	exampleContourA := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}},
		},
	}

	err = contourRepo.CreateContour(exampleContourA)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	exampleContourB := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{20, 20}, {30, 20}, {30, 30}, {20, 30}, {20, 20}}},
		},
	}
	err = contourRepo.CreateContour(exampleContourB)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	p.Suite.T().Run("GetPointsByContourID", func(t *testing.T) {
		tests := []struct {
			name           string
			countourID     uint
			found          bool
			expectedResult []models.Point
		}{
			{
				name:           "PointAExists",
				countourID:     exampleContourA.ID,
				found:          true,
				expectedResult: []models.Point{*examplePoint},
			},
			{
				name:       "PointBNotExists",
				countourID: exampleContourB.ID,
				found:      false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				points, err := repo.GetPointsByContourID(tt.countourID)
				if tt.found {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedResult, points)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, 0, len(points))
				}
			})
		}
	})

	tx.Rollback()
}
