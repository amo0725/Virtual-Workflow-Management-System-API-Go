package routes

import (
	"virtual_workflow_management_system_gin/controllers"
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/middlewares"

	"github.com/gin-gonic/gin"
)

func InitWorkflowRouter(routerGroup *gin.RouterGroup, resource *databases.Resource) {
	workflowController := controllers.NewWorkflowController(resource)

	authorizedGroup := routerGroup.Group("/workflows")
	authorizedGroup.Use(middlewares.JWTAuthMiddleware(resource.Redis))
	authorizedGroup.GET("", workflowController.GetWorkflows)
	authorizedGroup.GET("/:id", workflowController.GetWorkflow)
	authorizedGroup.POST("", workflowController.CreateWorkflow)
	authorizedGroup.PUT("/:id", workflowController.EditWorkflow)
	authorizedGroup.DELETE("/:id", workflowController.DeleteWorkflow)
	authorizedGroup.PUT("/:id/transfer/:username", workflowController.TransferWorkflow)
	authorizedGroup.GET("/:id/tasks", workflowController.GetTasks)
	authorizedGroup.GET("/:id/tasks/:taskID", workflowController.GetTask)
	authorizedGroup.POST("/:id/tasks", workflowController.CreateTask)
	authorizedGroup.PUT("/:id/tasks/:taskID", workflowController.EditTask)
}
