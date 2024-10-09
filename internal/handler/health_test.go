package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	// Initialize Gin
	gin.SetMode(gin.TestMode)

	// Define the test cases
	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Health check returns OK",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"status":"ok"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin engine instance
			router := gin.Default()

			// Register the Health handler to the engine
			router.GET("/health", Health)

			// Create a new HTTP request to the /health endpoint
			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Check if the status code matches the expected status code
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Check if the response body matches the expected body
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
