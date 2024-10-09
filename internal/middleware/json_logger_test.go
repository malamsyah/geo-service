package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestJSONLoggerMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Capture the logs output
	var buffer bytes.Buffer
	gin.DefaultWriter = io.MultiWriter(&buffer)

	// Create a new router and apply the middleware
	router := gin.New()
	router.Use(JSONLoggerMiddleware())

	// Define a sample route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Wait briefly to ensure log is written
	time.Sleep(10 * time.Millisecond)

	// Parse the log output
	logOutput := buffer.String()
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(logOutput), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Verify the log contains expected fields
	expectedFields := []string{
		"status_code",
		"path",
		"method",
		"start_time",
		"remote_addr",
		"response_time",
	}

	for _, field := range expectedFields {
		if _, exists := logEntry[field]; !exists {
			t.Errorf("Missing field '%s' in log entry", field)
		}
	}

	// Validate specific field values
	if logEntry["path"] != "/test" {
		t.Errorf("Expected path '/test', got '%v'", logEntry["path"])
	}
	if logEntry["method"] != "GET" {
		t.Errorf("Expected method 'GET', got '%v'", logEntry["method"])
	}
	if logEntry["status_code"] != float64(http.StatusOK) {
		t.Errorf("Expected status code %d, got %v", http.StatusOK, logEntry["status_code"])
	}
	if logEntry["remote_addr"] != "127.0.0.1" {
		t.Errorf("Expected remote_addr '127.0.0.1', got '%v'", logEntry["remote_addr"])
	}
}
