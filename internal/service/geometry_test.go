package service

import (
	"testing"

	"github.com/malamsyah/geo-service/internal/constants"
	"github.com/malamsyah/geo-service/internal/models"
	"github.com/malamsyah/geo-service/mocks/mock_internal/mock_repository"
	"go.uber.org/mock/gomock"
)

func TestGeometryService_IsValidPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
	svc := NewGeometryService(mockPointRepo, nil)
	tests := []struct {
		name           string
		point          *models.Point
		expectedResult bool
	}{
		{
			name: "ValidPoint",
			point: &models.Point{Data: models.Geometry{
				Type:             models.PointType,
				PointCoordinates: [2]float64{125.6, 10.1},
			}},
			expectedResult: true,
		},
		{
			name: "InvalidPoint",
			point: &models.Point{Data: models.Geometry{
				Type:             models.PointType,
				PointCoordinates: [2]float64{125.6, 200.1},
			}},
			expectedResult: false,
		},
		{
			name: "InvalidPointType",
			point: &models.Point{Data: models.Geometry{
				Type:             "RandomType",
				PointCoordinates: [2]float64{125.6, 200.1},
			}},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.IsValidPoint(tt.point); got != tt.expectedResult {
				t.Errorf("GeometryService.IsValidPoint() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

func TestGeometryService_CreatePoint(t *testing.T) {
	tests := []struct {
		name    string
		point   *models.Point
		mocks   func() *mock_repository.MockPointRepository
		wantErr bool
	}{
		{
			name: "ValidPoint",
			point: &models.Point{Data: models.Geometry{
				Type:             models.PointType,
				PointCoordinates: [2]float64{125.6, 10.1},
			}},
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().CreatePoint(gomock.Any()).Return(nil).Times(1)
				return mockPointRepo
			},
			wantErr: false,
		},
		{
			name: "InvalidPoint",
			point: &models.Point{Data: models.Geometry{
				Type:             models.PointType,
				PointCoordinates: [2]float64{180, 200.1},
			}},
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				return mockPointRepo
			},
			wantErr: true,
		},
		{
			name: "InvalidPointType",
			point: &models.Point{Data: models.Geometry{
				Type:             "RandomType",
				PointCoordinates: [2]float64{180, 180},
			}},
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				return mockPointRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := tt.mocks()
			svc := NewGeometryService(mockPointRepo, nil)

			if err := svc.CreatePoint(tt.point); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.CreatePoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_GetPoints(t *testing.T) {
	tests := []struct {
		name    string
		offset  int
		limit   int
		mocks   func() *mock_repository.MockPointRepository
		wantErr bool
	}{
		{
			name:   "OnePoint",
			offset: 0,
			limit:  1,
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().GetPoints(0, 1).Return([]models.Point{{ID: 1}}, nil).Times(1)
				return mockPointRepo
			},
			wantErr: false,
		},
		{
			name:   "ZeroLimit",
			offset: 0,
			limit:  0,
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().GetPoints(0, 0).Return([]models.Point{}, nil).Times(1)
				return mockPointRepo
			},
			wantErr: false,
		},
		{
			name:   "Error",
			offset: 0,
			limit:  10,
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().GetPoints(0, 10).Return(nil, constants.ErrInternal).Times(1)
				return mockPointRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := tt.mocks()
			svc := NewGeometryService(mockPointRepo, nil)

			if _, err := svc.GetPoints(tt.offset, tt.limit); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.GetPoints() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_GetPointByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mocks   func() *mock_repository.MockPointRepository
		wantErr bool
	}{
		{
			name: "ValidID",
			id:   1,
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().GetPointByID(uint(1)).Return(&models.Point{ID: 1}, nil).Times(1)
				return mockPointRepo
			},
			wantErr: false,
		},
		{
			name: "Error",
			id:   1,
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().GetPointByID(uint(1)).Return(nil, constants.ErrInternal).Times(1)
				return mockPointRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := tt.mocks()
			svc := NewGeometryService(mockPointRepo, nil)

			if _, err := svc.GetPointByID(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.GetPointByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_IsValidContour(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
	svc := NewGeometryService(nil, mockContourRepo)
	tests := []struct {
		name           string
		Contour        *models.Contour
		expectedResult bool
	}{
		{
			name: "ValidContour",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               models.PolygonType,
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
			}},
			expectedResult: true,
		},
		{
			name: "InvalidContour",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               models.PolygonType,
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.2}}},
			}},
			expectedResult: false,
		},
		{
			name: "InvalidContourType",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               "RandomType",
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
			}},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.IsValidContour(tt.Contour); got != tt.expectedResult {
				t.Errorf("GeometryService.IsValidContour() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

func TestGeometryService_CreateContour(t *testing.T) {
	tests := []struct {
		name    string
		Contour *models.Contour
		mocks   func() *mock_repository.MockContourRepository
		wantErr bool
	}{
		{
			name: "ValidContour",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               models.PolygonType,
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
			}},
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().CreateContour(gomock.Any()).Return(nil).Times(1)
				return mockContourRepo
			},
			wantErr: false,
		},
		{
			name: "InvalidContour",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               models.PolygonType,
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.2}}},
			}},
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				return mockContourRepo
			},
			wantErr: true,
		},
		{
			name: "InvalidContourType",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               "RandomType",
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.2}}},
			}},
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				return mockContourRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContourRepo := tt.mocks()
			svc := NewGeometryService(nil, mockContourRepo)

			if err := svc.CreateContour(tt.Contour); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.CreateContour() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_GetContours(t *testing.T) {
	tests := []struct {
		name    string
		offset  int
		limit   int
		mocks   func() *mock_repository.MockContourRepository
		wantErr bool
	}{
		{
			name:   "OneContour",
			offset: 0,
			limit:  1,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContours(0, 1).Return([]models.Contour{{ID: 1}}, nil).Times(1)
				return mockContourRepo
			},
			wantErr: false,
		},
		{
			name:   "ZeroLimit",
			offset: 0,
			limit:  0,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContours(0, 0).Return([]models.Contour{}, nil).Times(1)
				return mockContourRepo
			},
			wantErr: false,
		},
		{
			name:   "Error",
			offset: 0,
			limit:  10,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContours(0, 10).Return(nil, constants.ErrInternal).Times(1)
				return mockContourRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContourRepo := tt.mocks()
			svc := NewGeometryService(nil, mockContourRepo)

			if _, err := svc.GetContours(tt.offset, tt.limit); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.GetContours() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_GetContourByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mocks   func() *mock_repository.MockContourRepository
		wantErr bool
	}{
		{
			name: "ValidID",
			id:   1,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{ID: 1}, nil).Times(1)
				return mockContourRepo
			},
			wantErr: false,
		},
		{
			name: "Error",
			id:   1,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContourByID(uint(1)).Return(nil, constants.ErrInternal).Times(1)
				return mockContourRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContourRepo := tt.mocks()
			svc := NewGeometryService(nil, mockContourRepo)

			if _, err := svc.GetContourByID(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.GetContourByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_UpdateContour(t *testing.T) {
	tests := []struct {
		name    string
		Contour *models.Contour
		mocks   func() *mock_repository.MockContourRepository
		wantErr bool
	}{
		{
			name: "ValidContour",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               models.PolygonType,
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}},
			}},
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().UpdateContour(gomock.Any()).Return(nil).Times(1)
				return mockContourRepo
			},
			wantErr: false,
		},
		{
			name: "InvalidContour",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               models.PolygonType,
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.2}}},
			}},
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				return mockContourRepo
			},
			wantErr: true,
		},
		{
			name: "InvalidContourType",
			Contour: &models.Contour{Data: models.Geometry{
				Type:               "RandomType",
				PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.2}}},
			}},
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				return mockContourRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContourRepo := tt.mocks()
			svc := NewGeometryService(nil, mockContourRepo)

			if err := svc.UpdateContour(tt.Contour); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.UpdateContour() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_DeleteContour(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		mocks   func() *mock_repository.MockContourRepository
		wantErr bool
	}{
		{
			name: "ValidID",
			id:   1,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().DeleteContour(uint(1)).Return(nil).Times(1)
				return mockContourRepo
			},
			wantErr: false,
		},
		{
			name: "Error",
			id:   1,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().DeleteContour(uint(1)).Return(constants.ErrInternal).Times(1)
				return mockContourRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContourRepo := tt.mocks()
			svc := NewGeometryService(nil, mockContourRepo)

			if err := svc.DeleteContour(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.DeleteContour() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_GetPointsByContourID(t *testing.T) {
	tests := []struct {
		name      string
		contourID uint
		mocks     func() *mock_repository.MockPointRepository
		wantErr   bool
	}{
		{
			name:      "ValidID",
			contourID: 1,
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().GetPointsByContourID(uint(1)).Return([]models.Point{{ID: 1}}, nil).Times(1)
				return mockPointRepo
			},
			wantErr: false,
		},
		{
			name:      "Error",
			contourID: 1,
			mocks: func() *mock_repository.MockPointRepository {
				ctrl := gomock.NewController(t)
				mockPointRepo := mock_repository.NewMockPointRepository(ctrl)
				mockPointRepo.EXPECT().GetPointsByContourID(uint(1)).Return(nil, constants.ErrInternal).Times(1)
				return mockPointRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := tt.mocks()
			svc := NewGeometryService(mockPointRepo, nil)

			if _, err := svc.GetPointsByContourID(tt.contourID); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.GetPointsByContourID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeometryService_GetContoursIntersectArea(t *testing.T) {
	tests := []struct {
		name        string
		contourIDA  uint
		contourIDB  uint
		mocks       func() *mock_repository.MockContourRepository
		wantErr     bool
		expectedLen int
	}{
		{
			name:       "Valid",
			contourIDA: 1,
			contourIDB: 2,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{ID: 1}, nil).Times(1)
				mockContourRepo.EXPECT().GetContourByID(uint(2)).Return(&models.Contour{ID: 2}, nil).Times(1)
				mockContourRepo.EXPECT().GetContoursIntersectArea(uint(1), uint(2)).Return([]models.Contour{
					{ID: 1, Data: models.Geometry{Type: models.PolygonType, PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}}}},
					{ID: 2, Data: models.Geometry{Type: models.PolygonType, PolygonCoordinates: [][][2]float64{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}}}}}, nil).Times(1)
				return mockContourRepo
			},
			wantErr:     false,
			expectedLen: 2,
		},
		{
			name:       "Valid with a Multipolygon",
			contourIDA: 1,
			contourIDB: 2,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{ID: 1}, nil).Times(1)
				mockContourRepo.EXPECT().GetContourByID(uint(2)).Return(&models.Contour{ID: 2}, nil).Times(1)
				mockContourRepo.EXPECT().GetContoursIntersectArea(uint(1), uint(2)).Return([]models.Contour{{ID: 1, Data: models.Geometry{
					Type:                    models.MultiPolygon,
					MultiPolygonCoordinates: [][][][2]float64{{{{125.6, 10.1}, {125.7, 10.2}, {125.8, 10.3}, {125.6, 10.1}}}},
				}}}, nil).Times(1)
				return mockContourRepo
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name:       "InvalidContourA",
			contourIDA: 1,
			contourIDB: 2,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContourByID(uint(1)).Return(nil, constants.ErrNotFound).Times(1)
				return mockContourRepo
			},
			wantErr:     true,
			expectedLen: 0,
		},
		{
			name:       "InvalidContourB",
			contourIDA: 1,
			contourIDB: 2,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{ID: 1}, nil).Times(1)
				mockContourRepo.EXPECT().GetContourByID(uint(2)).Return(nil, constants.ErrNotFound).Times(1)
				return mockContourRepo
			},
			wantErr:     true,
			expectedLen: 0,
		},
		{
			name:       "Error",
			contourIDA: 1,
			contourIDB: 2,
			mocks: func() *mock_repository.MockContourRepository {
				ctrl := gomock.NewController(t)
				mockContourRepo := mock_repository.NewMockContourRepository(ctrl)
				mockContourRepo.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{ID: 1}, nil).Times(1)
				mockContourRepo.EXPECT().GetContourByID(uint(2)).Return(&models.Contour{ID: 2}, nil).Times(1)
				mockContourRepo.EXPECT().GetContoursIntersectArea(uint(1), uint(2)).Return(nil, constants.ErrInternal).Times(1)
				return mockContourRepo
			},
			wantErr:     true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockContourRepo := tt.mocks()
			svc := NewGeometryService(nil, mockContourRepo)

			if _, err := svc.GetContoursIntersectArea(tt.contourIDA, tt.contourIDB); (err != nil) != tt.wantErr {
				t.Errorf("GeometryService.GetContoursIntersectArea() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
