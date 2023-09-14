package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/requests"
	"virtual_workflow_management_system_gin/responses"
	"virtual_workflow_management_system_gin/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	OKStatus     = "OK"
	InvalidInput = "Invalid input"
	HTTPStatusOK = http.StatusOK
)

type MockUserService struct {
	LogoutError             error
	GetUsersByUsernameError error
}

var _ services.IUserService = &MockUserService{}

func (m *MockUserService) Register(c *gin.Context, req requests.RegisterRequest) {
	responses.Ok(c)
}

func (m *MockUserService) Login(c *gin.Context, req requests.LoginRequest) {
	responses.OkWithData(c, gin.H{
		"access_token":  "access_token",
		"refresh_token": "refresh_token",
	})
}

func (m *MockUserService) RefreshToken(c *gin.Context, req requests.RefreshTokenRequest) {
	responses.OkWithData(c, gin.H{
		"access_token":  "access_token",
		"refresh_token": "refresh_token",
	})
}

func (m *MockUserService) Logout(username string) error {
	if m.LogoutError != nil {
		return m.LogoutError
	}
	return nil
}

func (m *MockUserService) GetUsersByUsername(username string) (*models.User, error) {
	if m.GetUsersByUsernameError != nil {
		return nil, m.GetUsersByUsernameError
	}
	return &models.User{}, nil
}

var (
	mockUserService = new(MockUserService)
	userController  = UserController{UserService: mockUserService}
	router          = gin.Default()
)

func init() {
	gin.SetMode(gin.TestMode)
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/refresh-token", userController.RefreshToken)
	router.POST("/logout", userController.Logout)
}

func performRequest(method, path string, body []byte) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
	router.ServeHTTP(w, c.Request)
	return w
}

func runTestCase(t *testing.T, method, path string, input interface{}, expectedStatus int, shouldFind string) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	w := performRequest(method, path, requestBody)

	if err != nil {
		t.Fatalf("Failed to create new request: %v", err)
	}

	assert.Equal(t, expectedStatus, w.Code)
	assert.Contains(t, w.Body.String(), shouldFind)
}

func TestNewUserController(t *testing.T) {
	mockResource := &databases.Resource{}
	controller := NewUserController(mockResource)

	assert.NotNil(t, controller)
	assert.NotNil(t, controller.UserService)
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.RegisterRequest
		expected int
		message  string
	}{
		{"Valid input", requests.RegisterRequest{Username: "test", Password: "test123456", Role: "Admin"}, HTTPStatusOK, OKStatus},
		{"Missing Username", requests.RegisterRequest{Username: "", Password: "test123456", Role: "Admin"}, HTTPStatusOK, InvalidInput},
		{"Missing Password", requests.RegisterRequest{Username: "test", Password: "", Role: "Admin"}, HTTPStatusOK, InvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, http.MethodPost, "/register", tt.input, tt.expected, tt.message)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.LoginRequest
		expected int
		message  string
	}{
		{"Valid input", requests.LoginRequest{Username: "test", Password: "test123456"}, HTTPStatusOK, OKStatus},
		{"Missing Username", requests.LoginRequest{Password: "test123456"}, HTTPStatusOK, InvalidInput},
		{"Missing Password", requests.LoginRequest{Username: "test"}, HTTPStatusOK, InvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, http.MethodPost, "/login", tt.input, tt.expected, tt.message)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.RefreshTokenRequest
		expected int
		message  string
	}{
		{"Valid input", requests.RefreshTokenRequest{RefreshToken: "refresh_token"}, HTTPStatusOK, OKStatus},
		{"Missing RefreshToken", requests.RefreshTokenRequest{}, HTTPStatusOK, InvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTestCase(t, http.MethodPost, "/refresh-token", tt.input, tt.expected, tt.message)
		})
	}
}
func TestLogout(t *testing.T) {
	t.Run("Successful Logout", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user", models.JWTUser{Username: "testUser"})

		// Use the real UserController or a mock without errors here
		userController.Logout(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "OK")
	})

	// Sub-test for logout with error
	t.Run("Logout With Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user", models.JWTUser{Username: "testUser"})

		mockUserService := &MockUserService{LogoutError: errors.New("failed to logout")}
		userController := UserController{UserService: mockUserService}

		userController.Logout(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to logout")
	})
}
