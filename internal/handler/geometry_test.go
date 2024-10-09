package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/malamsyah/geo-service/internal/constants"
	"github.com/malamsyah/geo-service/internal/models"
	"github.com/malamsyah/geo-service/mocks/mock_internal/mock_service"
	"go.uber.org/mock/gomock"
)

func TestCreatePoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestBody          string
	}{
		{
			name:                 "Create point returns Created",
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":0,"data":{"type":"Point","coordinates":[5.123456,10.123456]}}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().CreatePoint(gomock.Any()).Return(nil)
				return mock
			},
			requestBody: `{"data":{"type":"Point","coordinates":[5.123456,10.123456]}}`,
		},
		{
			name:                 "Create point returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestBody: `{"data":{"type":"Point","coordinates":[5.123456,10.123456]}`,
		},
		{
			name:                 "Create point returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().CreatePoint(gomock.Any()).Return(constants.ErrInternal)
				return mock
			},
			requestBody: `{"data":{"type":"Point","coordinates":[5.123456,10.123456]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodPost, "/points", strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestGetPoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestParams        string
	}{
		{
			name:                 "Get points returns OK",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":0,"next":"http://localhost/points?page=1","previous":null,"results":[]}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetPoints(0, 10).Return([]models.Point{}, nil)
				return mock
			},
			requestParams: "page=0",
		},
		{
			name:                 "Get points returns OK with data",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":1,"next":"http://localhost/points?page=2","previous":"http://localhost/points?page=0","results":[{"id":1,"data":{"type":"Point","coordinates":[5.123456,10.123456]}}]}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetPoints(10, 10).Return([]models.Point{
					{
						ID: 1,
						Data: models.Geometry{
							Type:             "Point",
							PointCoordinates: [2]float64{5.123456, 10.123456},
						},
					},
				}, nil)
				return mock
			},
			requestParams: "page=1",
		},
		{
			name:                 "Get points returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestParams: "page=a",
		},
		{
			name:                 "Get points returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetPoints(0, 10).Return(nil, constants.ErrInternal)
				return mock
			},
			requestParams: "page=0",
		},
		{
			name:                 "Get points with contour ID returns OK",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":0,"next":"http://localhost/points?page=1","previous":null,"results":[]}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetPointsByContourID(uint(1)).Return([]models.Point{}, nil)
				return mock
			},
			requestParams: "contour=1",
		},
		{
			name:                 "Get points with contour ID returns OK with data",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":1,"next":"http://localhost/points?page=1","previous":null,"results":[{"id":1,"data":{"type":"Point","coordinates":[5.123456,10.123456]}}]}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetPointsByContourID(uint(1)).Return([]models.Point{
					{
						ID: 1,
						Data: models.Geometry{
							Type:             "Point",
							PointCoordinates: [2]float64{5.123456, 10.123456},
						},
					},
				}, nil)
				return mock
			},
			requestParams: "contour=1",
		},
		{
			name:                 "Get points with contour ID returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestParams: "contour=a",
		},
		{
			name:                 "Get points with contour ID returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetPointsByContourID(uint(1)).Return(nil, constants.ErrInternal)
				return mock
			},
			requestParams: "contour=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "http://localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodGet, "/points?"+tt.requestParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestCreateContour(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestBody          string
	}{
		{
			name:                 "Create Contour returns Created",
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().CreateContour(gomock.Any()).Return(nil)
				return mock
			},
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
		},
		{
			name:                 "Create Contour returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}`,
		},
		{
			name:                 "Create Contour returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().CreateContour(gomock.Any()).Return(constants.ErrInternal)
				return mock
			},
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodPost, "/contours", strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestGetContours(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestParams        string
	}{
		{
			name:                 "Get Contours returns OK",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":0,"next":"http://localhost/contours?page=1","previous":null,"results":[]}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContours(0, 10).Return([]models.Contour{}, nil)
				return mock
			},
			requestParams: "page=0",
		},
		{
			name:                 "Get Contours returns OK with data",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"count":1,"next":"http://localhost/contours?page=2","previous":"http://localhost/contours?page=0","results":[{"id":1,"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}]}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContours(10, 10).Return([]models.Contour{
					{
						ID: 1,
						Data: models.Geometry{
							Type:               "Polygon",
							PolygonCoordinates: [][][2]float64{{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}}},
						},
					},
				}, nil)
				return mock
			},
			requestParams: "page=1",
		},
		{
			name:                 "Get Contours returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestParams: "page=a",
		},
		{
			name:                 "Get Contours returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContours(0, 10).Return(nil, constants.ErrInternal)
				return mock
			},
			requestParams: "page=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "http://localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodGet, "/contours?"+tt.requestParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestGetContourByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestPath          string
	}{
		{
			name:                 "Get Contour by ID returns OK",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{
					ID: uint(1),
					Data: models.Geometry{
						Type:               "Polygon",
						PolygonCoordinates: [][][2]float64{{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}}},
					},
				}, nil).Times(1)
				return mock
			},
			requestPath: "/1",
		},
		{
			name:                 "Get Contour by ID returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestPath: "/a",
		},
		{
			name:                 "Get Contour by ID returns NotFound",
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"not found"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{}, constants.ErrNotFound)
				return mock
			},
			requestPath: "/1",
		},
		{
			name:                 "Get Contour by ID returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContourByID(uint(1)).Return(&models.Contour{}, constants.ErrInternal)
				return mock
			},
			requestPath: "/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "http://localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodGet, "/contours"+tt.requestPath, nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestUpdateContour(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestPath          string
		requestBody          string
	}{
		{
			name:                 "Update Contour returns OK",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1,"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().UpdateContour(gomock.Any()).Return(nil)
				return mock
			},
			requestPath: "/1",
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
		},
		{
			name:                 "Update Contour returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"unexpected EOF"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestPath: "/1",
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}`,
		},
		{
			name:                 "Update Contour returns BadRequest invalid params",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestPath: "/a",
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
		},
		{
			name:                 "Update Contour returns NotFound",
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"not found"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().UpdateContour(gomock.Any()).Return(constants.ErrNotFound)
				return mock
			},
			requestPath: "/1",
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
		},
		{
			name:                 "Update Contour returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().UpdateContour(gomock.Any()).Return(constants.ErrInternal)
				return mock
			},
			requestPath: "/1",
			requestBody: `{"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "http://localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodPut, "/contours"+tt.requestPath, strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestDeleteContour(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestPath          string
	}{
		{
			name:                 "Delete Contour returns NoContent",
			expectedStatusCode:   http.StatusNoContent,
			expectedResponseBody: "",
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().DeleteContour(uint(1)).Return(nil)
				return mock
			},
			requestPath: "/1",
		},
		{
			name:                 "Delete Contour returns BadRequest",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestPath: "/a",
		},
		{
			name:                 "Delete Contour returns NotFound",
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"error":"not found"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().DeleteContour(uint(1)).Return(constants.ErrNotFound)
				return mock
			},
			requestPath: "/1",
		},
		{
			name:                 "Delete Contour returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().DeleteContour(uint(1)).Return(constants.ErrInternal)
				return mock
			},
			requestPath: "/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "http://localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodDelete, "/contours"+tt.requestPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mocks                func() *mock_service.MockGeometryService
		requestParams        string
	}{
		{
			name:                 "Intersect returns OK",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[]`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContoursIntersectArea(uint(1), uint(2)).Return([]models.Contour{}, nil)
				return mock
			},
			requestParams: "contour_1=1&contour_2=2",
		},
		{
			name:                 "Intersect returns OK with data",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"id":1,"data":{"type":"Polygon","coordinates":[[[30,10],[40,40],[20,40],[10,20],[30,10]]]}}]`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContoursIntersectArea(uint(1), uint(2)).Return([]models.Contour{
					{
						ID: 1,
						Data: models.Geometry{
							Type:               "Polygon",
							PolygonCoordinates: [][][2]float64{{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}}},
						},
					},
				}, nil)
				return mock
			},
			requestParams: "contour_1=1&contour_2=2",
		},
		{
			name:                 "Intersect returns BadRequest contour_1",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestParams: "contour_1=a&contour_2=2",
		},
		{
			name:                 "Intersect returns BadRequest contour_2",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"strconv.Atoi: parsing \"a\": invalid syntax"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				return mock
			},
			requestParams: "contour_1=1&contour_2=a",
		},
		{
			name:                 "Intersect returns InternalServerError",
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal error"}`,
			mocks: func() *mock_service.MockGeometryService {
				ctrl := gomock.NewController(t)
				mock := mock_service.NewMockGeometryService(ctrl)
				mock.EXPECT().GetContoursIntersectArea(uint(1), uint(2)).Return(nil, constants.ErrInternal)
				return mock
			},
			requestParams: "contour_1=1&contour_2=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			handler := NewGeometryHandler(tt.mocks(), "http://localhost")
			handler.RegisterRoutes(router.Group("/"))

			req, err := http.NewRequest(http.MethodGet, "/intersections?"+tt.requestParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatusCode, w.Code)
			}

			if w.Body.String() != tt.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", tt.expectedResponseBody, w.Body.String())
			}
		})
	}
}
