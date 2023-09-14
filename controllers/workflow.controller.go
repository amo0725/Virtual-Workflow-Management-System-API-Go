package controllers

import (
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/requests"
	"virtual_workflow_management_system_gin/responses"
	"virtual_workflow_management_system_gin/services"

	"github.com/gin-gonic/gin"
)

type WorkflowController struct {
	WorkflowService services.IWorkflowService
	UserService     services.IUserService
}

func NewWorkflowController(resource *databases.Resource) *WorkflowController {
	workflowService := services.NewWorkflowService(resource)
	userService := services.NewUserService(resource)
	return &WorkflowController{WorkflowService: workflowService, UserService: userService}
}

// @Security access_token
// @Summary Get all workflows
// @Tags Workflows
// @version 1.0
// @Description Get all workflows
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} string "OK"
// @Router /workflows [get]
func (controller *WorkflowController) GetWorkflows(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	workflows, err := controller.WorkflowService.GetWorkflows(user.Username)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"workflows": workflows,
	})
}

// @Security access_token
// @Summary Get a workflow
// @Tags Workflows
// @version 1.0
// @Description Get a workflow by ID
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Success 200 {object} string "OK"
// @Router /workflows/{id} [get]
func (controller *WorkflowController) GetWorkflow(c *gin.Context) {
	workflowID := c.Param("id")

	workflow, err := controller.WorkflowService.GetWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"workflow": workflow,
	})
}

// @Security access_token
// @Summary Create a new workflow
// @Tags Workflows
// @version 1.0
// @Description Create a new workflow with the input payload
// @Accept  application/json
// @Produce  application/json
// @Param workflow body requests.CreateWorkflowRequest true "Workflow for creation"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /workflows [post]
func (controller *WorkflowController) CreateWorkflow(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	var req requests.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, "Invalid input")
		return
	}

	insertedID, err := controller.WorkflowService.CreateWorkflow(user.Username, req)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"workflow_id": insertedID,
	})

}

// @Security access_token
// @Summary Edit an existing workflow
// @Tags Workflows
// @version 1.0
// @Description Edit an existing workflow with the input payload
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Param workflow body requests.EditWorkflowRequest true "Workflow for editing"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /workflows/{id} [put]
func (controller *WorkflowController) EditWorkflow(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	workflowID := c.Param("id")

	workflow, err := controller.WorkflowService.GetWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, "failed to get workflow")
		return
	}

	if !workflow.CheckWorkflowAccess(user, "delete") {
		responses.Error(c, "unauthorized")
		return
	}

	var req requests.EditWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, "Invalid input")
		return
	}
	updatedWorkflow, err := controller.WorkflowService.EditWorkflowByID(workflowID, req)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"workflow": updatedWorkflow,
	})
}

// @Security access_token
// @Summary Delete a workflow
// @Tags Workflows
// @version 1.0
// @Description Delete a workflow by ID
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /workflows/{id} [delete]
func (controller *WorkflowController) DeleteWorkflow(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	workflowID := c.Param("id")

	workflow, err := controller.WorkflowService.GetWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, "failed to get workflow")
		return
	}

	if !workflow.CheckWorkflowAccess(user, "delete") {
		responses.Error(c, "unauthorized")
		return
	}

	err = controller.WorkflowService.DeleteWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.Ok(c)
}

// @Security access_token
// @Summary Transfer a workflow
// @Tags Workflows
// @version 1.0
// @Description Transfer a workflow by ID
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Param username path string true "Username"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /workflows/{id}/transfer/{username} [put]
func (controller *WorkflowController) TransferWorkflow(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	workflowID := c.Param("id")
	newOwner := c.Param("username")

	workflow, err := controller.WorkflowService.GetWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, "failed to get workflow")
		return
	}

	if !workflow.CheckWorkflowAccess(user, "transfer") {
		responses.Error(c, "unauthorized")
		return
	}

	_, err = controller.UserService.GetUsersByUsername(newOwner)
	if err != nil {
		responses.Error(c, "user does not exist")
		return
	}

	transferedWorkflow, err := controller.WorkflowService.TransferWorkflowByID(workflowID, newOwner)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"workflow": transferedWorkflow,
	})
}

// @Security access_token
// @Summary Get all tasks
// @Tags Workflows
// @version 1.0
// @Description Get all tasks
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Success 200 {object} string "OK"
// @Router /workflows/{id}/tasks [get]
func (controller *WorkflowController) GetTasks(c *gin.Context) {
	workflowID := c.Param("id")

	tasks, err := controller.WorkflowService.GetTasksByWorkflowID(workflowID)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"tasks": tasks,
	})
}

// @Security access_token
// @Summary Get a task
// @Tags Workflows
// @version 1.0
// @Description Get a task by ID
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Param taskID path string true "Task ID"
// @Success 200 {object} string "OK"
// @Router /workflows/{id}/tasks/{taskID} [get]
func (controller *WorkflowController) GetTask(c *gin.Context) {
	workflowID := c.Param("id")
	taskID := c.Param("taskID")

	task, err := controller.WorkflowService.GetTaskByID(workflowID, taskID)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"task": task,
	})
}

// @Security access_token
// @Summary Create a task
// @Tags Workflows
// @version 1.0
// @Description Create a task with the input payload
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Param task body requests.CreateTaskRequest true "Task for creation"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /workflows/{id}/tasks [post]
func (controller *WorkflowController) CreateTask(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	workflowID := c.Param("id")

	workflow, err := controller.WorkflowService.GetWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, "failed to get workflow")
		return
	}

	if !workflow.CheckWorkflowAccess(user, "delete") {
		responses.Error(c, "unauthorized")
		return
	}

	var req requests.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, "Invalid input")
		return
	}
	createdTaskID, err := controller.WorkflowService.CreateTaskByWorkflowID(workflowID, req)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"task_id": createdTaskID,
	})
}

// @Security access_token
// @Summary Edit a task
// @Tags Workflows
// @version 1.0
// @Description Edit a task with the input payload
// @Accept  application/json
// @Produce  application/json
// @Param workflow_id path string true "Workflow ID"
// @Param task_id path string true "Task ID"
// @Param task body requests.EditTaskRequest true "Task for editing"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /workflows/{workflow_id}/tasks/{task_id} [put]
func (controller *WorkflowController) EditTask(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	workflowID := c.Param("id")
	taskID := c.Param("taskID")

	workflow, err := controller.WorkflowService.GetWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, "failed to get workflow")
		return
	}

	if !workflow.CheckWorkflowAccess(user, "delete") {
		responses.Error(c, "unauthorized")
		return
	}

	var req requests.EditTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, "Invalid input")
		return
	}
	task, err := controller.WorkflowService.EditTaskByID(workflowID, taskID, req)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.OkWithData(c, gin.H{
		"task": task,
	})
}

// @Security access_token
// @Summary Delete a task
// @Tags Workflows
// @version 1.0
// @Description Delete a task by ID
// @Accept  application/json
// @Produce  application/json
// @Param id path string true "Workflow ID"
// @Param taskID path string true "Task ID"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /workflows/{id}/tasks/{taskID} [delete]
func (controller *WorkflowController) DeleteTask(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	workflowID := c.Param("id")
	taskID := c.Param("taskID")

	workflow, err := controller.WorkflowService.GetWorkflowByID(workflowID)
	if err != nil {
		responses.Error(c, "failed to get workflow")
		return
	}

	if !workflow.CheckWorkflowAccess(user, "delete") {
		responses.Error(c, "unauthorized")
		return
	}

	err = controller.WorkflowService.DeleteTaskByID(workflowID, taskID)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.Ok(c)
}
