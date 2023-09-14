package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/requests"
	"virtual_workflow_management_system_gin/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockWorkflowService struct {
	GetWorkflowsError           error
	GetWorkflowByIDError        error
	CreateWorkflowError         error
	EditWorkflowByIDError       error
	DeleteWorkflowByIDError     error
	TransferWorkflowByIDError   error
	GetTasksByWorkflowIDError   error
	GetTaskByIDError            error
	CreateTaskByWorkflowIDError error
	EditTaskByIDError           error
	DeleteTaskByIDError         error
}

var _ services.IWorkflowService = &MockWorkflowService{}

func (m *MockWorkflowService) GetWorkflows(username string) ([]models.Workflow, error) {
	if m.GetWorkflowsError != nil {
		return nil, m.GetWorkflowsError
	}
	return []models.Workflow{}, nil
}

func (m *MockWorkflowService) GetWorkflowByID(workflowID string) (*models.Workflow, error) {
	if m.GetWorkflowByIDError != nil {
		return nil, m.GetWorkflowByIDError
	}
	return &models.Workflow{Owner: "testUser"}, nil
}

func (m *MockWorkflowService) CreateWorkflow(username string, req requests.CreateWorkflowRequest) (*string, error) {
	if m.CreateWorkflowError != nil {
		return nil, m.CreateWorkflowError
	}
	id := "newID"
	return &id, nil
}

func (m *MockWorkflowService) EditWorkflowByID(workflowID string, req requests.EditWorkflowRequest) (*models.Workflow, error) {
	if m.EditWorkflowByIDError != nil {
		return nil, m.EditWorkflowByIDError
	}
	return &models.Workflow{}, nil
}

func (m *MockWorkflowService) DeleteWorkflowByID(workflowID string) error {
	if m.DeleteWorkflowByIDError != nil {
		return m.DeleteWorkflowByIDError
	}
	return nil
}

func (m *MockWorkflowService) TransferWorkflowByID(workflowID string, username string) (*models.Workflow, error) {
	if m.TransferWorkflowByIDError != nil {
		return nil, m.TransferWorkflowByIDError
	}
	return &models.Workflow{}, nil
}

func (m *MockWorkflowService) GetTasksByWorkflowID(workflowID string) ([]models.Task, error) {
	if m.GetTasksByWorkflowIDError != nil {
		return nil, m.GetTasksByWorkflowIDError
	}
	return []models.Task{}, nil
}

func (m *MockWorkflowService) GetTaskByID(workflowID string, taskID string) (*models.Task, error) {
	if m.GetTaskByIDError != nil {
		return nil, m.GetTaskByIDError
	}
	return &models.Task{}, nil
}

func (m *MockWorkflowService) CreateTaskByWorkflowID(workflowID string, req requests.CreateTaskRequest) (*string, error) {
	if m.CreateTaskByWorkflowIDError != nil {
		return nil, m.CreateTaskByWorkflowIDError
	}
	id := "newID"
	return &id, nil
}

func (m *MockWorkflowService) EditTaskByID(workflowID string, taskID string, req requests.EditTaskRequest) (*models.Task, error) {
	if m.EditTaskByIDError != nil {
		return nil, m.EditTaskByIDError
	}
	return &models.Task{}, nil
}

func (m *MockWorkflowService) DeleteTaskByID(workflowID string, taskID string) error {
	if m.DeleteTaskByIDError != nil {
		return m.DeleteTaskByIDError
	}
	return nil
}

var (
	mockWorkflowService = new(MockWorkflowService)
	workflowController  = WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}
)

func init() {
	gin.SetMode(gin.TestMode)
	router.GET("/workflows", workflowController.GetWorkflows)
	router.GET("/workflows/:id", workflowController.GetWorkflow)
	router.POST("/workflows", workflowController.CreateWorkflow)
	router.PUT("/workflows/:id", workflowController.EditWorkflow)
	router.DELETE("/workflows/:id", workflowController.DeleteWorkflow)
	router.POST("/workflows/:id/transfer", workflowController.TransferWorkflow)
	router.GET("/workflows/:id/tasks", workflowController.GetTasks)
	router.GET("/workflows/:id/tasks/:taskID", workflowController.GetTask)
	router.POST("/workflows/:id/tasks", workflowController.CreateTask)
	router.PUT("/workflows/:id/tasks/:taskID", workflowController.EditTask)
	router.DELETE("/workflows/:id/tasks/:taskID", workflowController.DeleteTask)
}

func TestNewWorkflowController(t *testing.T) {
	mockResource := &databases.Resource{}
	controller := NewWorkflowController(mockResource)

	assert.NotNil(t, controller)
	assert.NotNil(t, controller.UserService)
}

func TestGetWorkflows(t *testing.T) {
	t.Run("Successful GetWorkflows", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows", nil)
		c.Set("user", models.JWTUser{Username: "testUser"})

		workflowController.GetWorkflows(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "OK")
	})

	t.Run("Failed GetWorkflows", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows", nil)
		c.Set("user", models.JWTUser{Username: "testUser"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowsError: errors.New("failed to retrieve workflows")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.GetWorkflows(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to retrieve workflows")
	})
}

func TestGetWorkflow(t *testing.T) {
	t.Run("Successful GetWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testUser"})

		workflowController.GetWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "OK")
	})

	t.Run("Failed GetWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testUser"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowByIDError: errors.New("workflow does not exist")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.GetWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "workflow does not exist")
	})
}

func TestCreateWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.CreateWorkflowRequest
		expected int
		message  string
	}{
		{"Valid input", requests.CreateWorkflowRequest{Name: "test"}, HTTPStatusOK, OKStatus},
		{"Missing RefreshToken", requests.CreateWorkflowRequest{}, HTTPStatusOK, InvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest(http.MethodPost, "/workflows", bytes.NewBuffer(requestBody))
			c.Set("user", models.JWTUser{Username: "testUser"})

			workflowController.CreateWorkflow(c)

			assert.Equal(t, tt.expected, w.Code)
			assert.Contains(t, w.Body.String(), tt.message)
		})
	}

	t.Run("Failed CreateWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/workflows", strings.NewReader(`{"name":"test"}`))
		c.Set("user", models.JWTUser{Username: "testUser"})

		mockWorkflowService := &MockWorkflowService{CreateWorkflowError: errors.New("failed to create workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.CreateWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code) // Update based on your error handling
		assert.Contains(t, w.Body.String(), "failed to create workflow")
	})
}

func TestEditWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.EditWorkflowRequest
		expected int
		message  string
	}{
		{"Valid input", requests.EditWorkflowRequest{Name: "edited"}, HTTPStatusOK, OKStatus},
		{"Invalid Input", requests.EditWorkflowRequest{}, HTTPStatusOK, InvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest(http.MethodPut, "/workflows/some_id", bytes.NewBuffer(requestBody))
			c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

			workflowController.EditWorkflow(c)

			assert.Equal(t, tt.expected, w.Code)
			assert.Contains(t, w.Body.String(), tt.message)
		})
	}

	t.Run("Failed EditWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/workflows/some_id", strings.NewReader(`{"name":"edited"}`))
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{EditWorkflowByIDError: errors.New("failed to edit workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.EditWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to edit workflow")
	})

	t.Run("Failed GetWorkflowByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/workflows/some_id", strings.NewReader(`{"name":"test"}`))
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowByIDError: errors.New("failed to get workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.EditWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get workflow")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/workflows/some_id", strings.NewReader(`{"name":"test"}`))
		c.Set("user", models.JWTUser{Username: "testWrongUser"})

		workflowController.EditWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

func TestDeleteWorkflow(t *testing.T) {

	t.Run("Successfull DeleteWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/workflows/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		workflowController.DeleteWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "OK")
	})

	t.Run("Failed DeleteWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/workflows/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{DeleteWorkflowByIDError: errors.New("failed to delete workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.DeleteWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to delete workflow")
	})

	t.Run("Failed GetWorkflowByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/workflows/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowByIDError: errors.New("failed to get workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.DeleteWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get workflow")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/workflows/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testWrongUser"})

		workflowController.DeleteWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

func TestTransferWorkflow(t *testing.T) {
	t.Run("Successful TransferWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/workflows/some_id/transfer/username", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		workflowController.TransferWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), OKStatus)
	})

	t.Run("Failed TransferWorkflow", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/workflows/some_id/transfer/username", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{TransferWorkflowByIDError: errors.New("failed to transfer workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.TransferWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to transfer workflow")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/workflows/some_id/transfer/username", nil)
		c.Set("user", models.JWTUser{Username: "testWrongUser", Role: "admin"})

		workflowController.TransferWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})

	t.Run("Failed GetWorkflowByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/workflows/some_id/transfer/username", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowByIDError: errors.New("failed to get workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.TransferWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get workflow")
	})

	t.Run("Failed GetUsersByUsername", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/workflows/some_id/transfer/username", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockUserService := &MockUserService{GetUsersByUsernameError: errors.New("user does not exist")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.TransferWorkflow(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user does not exist")
	})
}

func TestGetTasks(t *testing.T) {
	t.Run("Successful GetTasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows/some_id/tasks", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		workflowController.GetTasks(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), OKStatus)
	})

	t.Run("Failed GetTasks", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows/some_id/tasks", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetTasksByWorkflowIDError: errors.New("failed to get tasks")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.GetTasks(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get tasks")
	})
}

func TestGetTask(t *testing.T) {
	t.Run("Successful GetTask", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows/some_id/tasks/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		workflowController.GetTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), OKStatus)
	})

	t.Run("Failed GetTask", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/workflows/some_id/tasks/some_id", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetTaskByIDError: errors.New("failed to get task")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.GetTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get task")
	})
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.CreateTaskRequest
		expected int
		message  string
	}{
		{"Valid input", requests.CreateTaskRequest{Name: "test", Description: "description"}, HTTPStatusOK, OKStatus},
		{"Invalid Input", requests.CreateTaskRequest{Description: "description"}, HTTPStatusOK, InvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest(http.MethodPost, "/workflows/some_id/tasks", bytes.NewBuffer(requestBody))
			c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

			workflowController.CreateTask(c)

			assert.Equal(t, tt.expected, w.Code)
			assert.Contains(t, w.Body.String(), tt.message)
		})
	}

	t.Run("Failed GetWorkflowByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/workflows/some_id/tasks", nil)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowByIDError: errors.New("failed to get workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.CreateTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get workflow")
	})

	t.Run("Failed CreateTaskByWorkflowID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodPost,
			"/workflows/some_id/tasks",
			strings.NewReader(`{"name":"test","description":"description"}`),
		)
		c.Set("user", models.JWTUser{Username: "user", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{CreateTaskByWorkflowIDError: errors.New("failed to create task")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.CreateTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to create task")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodPost,
			"/workflows/some_id/tasks",
			strings.NewReader(`{"name":"test","description":"description"}`),
		)
		c.Set("user", models.JWTUser{Username: "testWrongUser"})

		workflowController.CreateTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

func TestEditTask(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.EditTaskRequest
		expected int
		message  string
	}{
		{"Valid input", requests.EditTaskRequest{Name: "test", Description: "description", Status: "Pending", Order: 1}, HTTPStatusOK, OKStatus},
		{"Invalid Input", requests.EditTaskRequest{}, HTTPStatusOK, InvalidInput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			requestBody, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest(http.MethodPut, "/workflows/some_id/tasks/some_id", bytes.NewBuffer(requestBody))
			c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

			workflowController.EditTask(c)

			assert.Equal(t, tt.expected, w.Code)
			assert.Contains(t, w.Body.String(), tt.message)
		})
	}

	t.Run("Failed EditTask", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodPut,
			"/workflows/some_id/tasks/some_id",
			strings.NewReader(`{"name":"test","description":"description","status":"Pending","order":1}`),
		)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{EditTaskByIDError: errors.New("failed to edit task")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.EditTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to edit task")
	})

	t.Run("Failed GetWorkflowByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodPut,
			"/workflows/some_id/tasks/some_id",
			strings.NewReader(`{"name":"test","description":"description","status":"Pending","order":1}`),
		)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowByIDError: errors.New("failed to get workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.EditTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get workflow")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodPut,
			"/workflows/some_id/tasks/some_id",
			strings.NewReader(`{"name":"test","description":"description","status":"Pending","order":1}`),
		)
		c.Set("user", models.JWTUser{Username: "testWrongUser"})

		workflowController.EditTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("Successful DeleteTask", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodDelete,
			"/workflows/some_id/tasks/some_id",
			nil,
		)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		workflowController.DeleteTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), OKStatus)
	})

	t.Run("Failed DeleteTask", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodDelete,
			"/workflows/some_id/tasks/some_id",
			nil,
		)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{DeleteTaskByIDError: errors.New("failed to delete task")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.DeleteTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to delete task")
	})

	t.Run("Failed GetWorkflowByID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodDelete,
			"/workflows/some_id/tasks/some_id",
			nil,
		)
		c.Set("user", models.JWTUser{Username: "testUser", Role: "admin"})

		mockWorkflowService := &MockWorkflowService{GetWorkflowByIDError: errors.New("failed to get workflow")}
		workflowController := WorkflowController{WorkflowService: mockWorkflowService, UserService: mockUserService}

		workflowController.DeleteTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "failed to get workflow")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(
			http.MethodDelete,
			"/workflows/some_id/tasks/some_id",
			nil,
		)
		c.Set("user", models.JWTUser{Username: "testWrongUser"})

		workflowController.DeleteTask(c)

		assert.Equal(t, HTTPStatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "unauthorized")
	})
}
