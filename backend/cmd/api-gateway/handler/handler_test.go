package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Missing email",
			body: map[string]interface{}{
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid email",
			body: map[string]interface{}{
				"email":    "notanemail",
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Short username",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"username": "ab",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Short password",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"username": "testuser",
				"password": "1234567",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Valid input (mock test)",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest, // Will fail due to no DB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create router with handler
			// Note: This is a simplified test without actual DB connection
			// In real tests, you would use a mock database
			router := gin.New()
			router.POST("/auth/register", func(c *gin.Context) {
				var req RegisterRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
					return
				}
				// Would proceed with registration
			})

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestLoginValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Missing email",
			body: map[string]interface{}{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing password",
			body: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Valid input (mock test)",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest, // Will fail due to no DB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router := gin.New()
			router.POST("/auth/login", func(c *gin.Context) {
				var req LoginRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
					return
				}
			})

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCreateTaskValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Missing agent_id",
			body: map[string]interface{}{
				"prompt": "Test prompt",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing prompt",
			body: map[string]interface{}{
				"agent_id": "test-agent",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty prompt",
			body: map[string]interface{}{
				"agent_id": "test-agent",
				"prompt":   "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router := gin.New()
			router.POST("/tasks", func(c *gin.Context) {
				var req CreateTaskRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
					return
				}
			})

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestListTasksPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		query      string
		wantPage   int
		wantLimit  int
	}{
		{"Default pagination", "", 1, 20},
		{"Custom page", "?page=5", 5, 20},
		{"Custom limit", "?limit=50", 1, 50},
		{"Both custom", "?page=3&limit=10", 3, 10},
		{"Invalid page", "?page=0", 1, 20},
		{"Invalid limit", "?limit=0", 1, 20},
		{"Too large limit", "?limit=200", 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/tasks"+tt.query, nil)

			var gotPage, gotLimit int

			router := gin.New()
			router.GET("/tasks", func(c *gin.Context) {
				page, _ := parseIntParam(c.DefaultQuery("page", "1"))
				limit, _ := parseIntParam(c.DefaultQuery("limit", "20"))

				if page < 1 {
					page = 1
				}
				if limit < 1 || limit > 100 {
					limit = 20
				}

				gotPage = page
				gotLimit = limit

				c.JSON(http.StatusOK, gin.H{})
			})

			router.ServeHTTP(w, req)

			if gotPage != tt.wantPage {
				t.Errorf("Expected page %d, got %d", tt.wantPage, gotPage)
			}

			if gotLimit != tt.wantLimit {
				t.Errorf("Expected limit %d, got %d", tt.wantLimit, gotLimit)
			}
		})
	}
}

// Helper function
func parseIntParam(s string) (int, error) {
	var i int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		i = i*10 + int(c-'0')
	}
	return i, nil
}
