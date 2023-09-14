package requests

import "virtual_workflow_management_system_gin/models"

type CreateWorkflowRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

type CreateTaskRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description"`
}

type EditWorkflowRequest struct {
	Name string `json:"name" binding:"min=3,max=100"`
}

type EditTaskRequest struct {
	Name        string            `json:"name" binding:"required,min=3,max=100"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status" binding:"required"`
	Order       int               `json:"order" binding:"required"`
}
