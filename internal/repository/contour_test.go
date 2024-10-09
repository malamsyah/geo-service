package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/malamsyah/geo-service/internal/models"
)

type ContourRepoTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func TestContourRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ContourRepoTestSuite))
}

func (p *ContourRepoTestSuite) SetupSuite() {
	p.db = setupTestDB(p.Suite.T())
}

func (p *ContourRepoTestSuite) TestContourRepository_CreateContour() {
	tx := p.db.Begin()
	repo := NewContourRepository(tx)
	p.Suite.T().Run("CreateContour", func(t *testing.T) {
		tests := []struct {
			name    string
			Contour *models.Contour
			wantErr bool
		}{
			{
				name: "ValidContour",
				Contour: &models.Contour{
					Data: models.Geometry{
						Type:               "Polygon",
						PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
					},
				},
				wantErr: false,
			},
			{
				name: "EmptyData",
				Contour: &models.Contour{
					Data: models.Geometry{
						Type: "Polygon",
					},
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.CreateContour(tt.Contour)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, tt.Contour.ID)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *ContourRepoTestSuite) TestContourRepository_GetContourByID() {
	tx := p.db.Begin()
	repo := NewContourRepository(tx)

	exampleContour := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
		},
	}
	err := repo.CreateContour(exampleContour)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	p.Suite.T().Run("GetContourByID", func(t *testing.T) {
		tests := []struct {
			name    string
			ID      uint
			wantErr bool
		}{
			{
				name:    "ContourExists",
				ID:      exampleContour.ID,
				wantErr: false,
			},
			{
				name:    "ContourNotExists",
				ID:      999999,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				actualContour, err := repo.GetContourByID(tt.ID)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, actualContour.ID)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *ContourRepoTestSuite) TestContourRepository_GetContours() {
	tx := p.db.Begin()
	repo := NewContourRepository(tx)

	exampleContour := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
		},
	}
	err := repo.CreateContour(exampleContour)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	p.Suite.T().Run("GetContours", func(t *testing.T) {
		tests := []struct {
			name           string
			limit          int
			offset         int
			expectedLength int
			wantErr        bool
		}{
			{
				name:           "OneContour",
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
				Contours, err := repo.GetContours(tt.offset, tt.limit)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedLength, len(Contours))
				}
			})
		}
	})

	tx.Rollback()
}

func (p *ContourRepoTestSuite) TestContourRepository_UpdateContour() {
	tx := p.db.Begin()
	repo := NewContourRepository(tx)

	exampleContour := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
		},
	}
	err := repo.CreateContour(exampleContour)
	if err != nil {
		p.Suite.T().Fatal(err)
	}
	p.Suite.T().Run("UpdateContour", func(t *testing.T) {
		tests := []struct {
			name    string
			Contour *models.Contour
			wantErr bool
		}{
			{
				name: "ValidContour",
				Contour: &models.Contour{
					ID: exampleContour.ID,
					Data: models.Geometry{
						Type:               "Polygon",
						PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
					},
				},
				wantErr: false,
			},
			{
				name: "EmptyData",
				Contour: &models.Contour{
					ID: exampleContour.ID,
					Data: models.Geometry{
						Type:               "Polygon",
						PolygonCoordinates: [][][2]float64{},
					},
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.UpdateContour(tt.Contour)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, tt.Contour.ID)
					assert.Equal(t, exampleContour.ID, tt.Contour.ID)

					actualContour, err := repo.GetContourByID(tt.Contour.ID)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, tt.Contour.Data, actualContour.Data)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *ContourRepoTestSuite) TestContourRepository_DeleteContour() {
	tx := p.db.Begin()
	repo := NewContourRepository(tx)

	exampleContour := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
		},
	}
	err := repo.CreateContour(exampleContour)
	if err != nil {
		p.Suite.T().Fatal(err)
	}
	p.Suite.T().Run("DeleteContour", func(t *testing.T) {
		tests := []struct {
			name    string
			ID      uint
			wantErr bool
		}{
			{
				name:    "ContourExists",
				ID:      exampleContour.ID,
				wantErr: false,
			},
			{
				name:    "ContourNotExists",
				ID:      999999,
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.DeleteContour(tt.ID)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					actualContour, err := repo.GetContourByID(tt.ID)
					assert.Error(t, err)
					assert.Nil(t, actualContour)
				}
			})
		}
	})

	tx.Rollback()
}

func (p *ContourRepoTestSuite) TestContourRepository_GetContoursIntersectArea() {
	tx := p.db.Begin()
	repo := NewContourRepository(tx)

	exampleContourA := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
		},
	}
	err := repo.CreateContour(exampleContourA)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	exampleContourB := &models.Contour{
		Data: models.Geometry{
			Type:               "Polygon",
			PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
		},
	}
	err = repo.CreateContour(exampleContourB)
	if err != nil {
		p.Suite.T().Fatal(err)
	}

	p.Suite.T().Run("GetContoursIntersectArea", func(t *testing.T) {
		tests := []struct {
			name  string
			IDA   uint
			IDB   uint
			found bool
		}{
			{
				name:  "ContourExists",
				IDA:   exampleContourA.ID,
				IDB:   exampleContourB.ID,
				found: true,
			},
			{
				name:  "ContourNotExists",
				IDA:   999999,
				IDB:   999999,
				found: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				Contours, err := repo.GetContoursIntersectArea(tt.IDA, tt.IDB)
				if !tt.found {
					assert.NoError(t, err)
					assert.Zero(t, len(Contours))
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, len(Contours))
				}
			})
		}
	})

	tx.Rollback()
}
