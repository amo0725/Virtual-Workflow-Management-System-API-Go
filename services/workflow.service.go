package services

import (
	"context"
	"virtual_workflow_management_system_gin/common"
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/repositories"
	"virtual_workflow_management_system_gin/requests"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var WorkflowService IWorkflowService

type workflowService struct {
	workflowEntity repositories.IWorkflow
	mongoClient    *mongo.Client
}

type IWorkflowService interface {
	GetWorkflows(username string) ([]models.Workflow, error)
	GetWorkflowByID(workflowID string) (*models.Workflow, error)
	CreateWorkflow(username string, req requests.CreateWorkflowRequest) (*string, error)
	EditWorkflowByID(workflowID string, req requests.EditWorkflowRequest) (*models.Workflow, error)
	DeleteWorkflowByID(workflowID string) error
	TransferWorkflowByID(workflowID string, username string) (*models.Workflow, error)
	GetTasksByWorkflowID(workflowID string) ([]models.Task, error)
	GetTaskByID(workflowID string, taskID string) (*models.Task, error)
	CreateTaskByWorkflowID(workflowID string, req requests.CreateTaskRequest) (*string, error)
	EditTaskByID(workflowID string, taskID string, req requests.EditTaskRequest) (*models.Task, error)
	DeleteTaskByID(workflowID string, taskID string) error
}

func NewWorkflowService(resource *databases.Resource) *workflowService {
	if resource == nil || resource.MongoDB == nil || resource.Redis == nil {
		return &workflowService{}
	}
	return &workflowService{
		workflowEntity: repositories.NewWorkflowEntity(resource),
		mongoClient:    resource.MongoDB.Client(),
	}
}

func (service *workflowService) GetWorkflows(username string) ([]models.Workflow, error) {
	workflows, err := service.workflowEntity.FindWorkflowsByUsername(username)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return workflows, nil
}

func (service *workflowService) GetWorkflowByID(workflowID string) (*models.Workflow, error) {
	workflow, err := service.workflowEntity.FindWorkflowByID(workflowID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return workflow, nil
}

func (service *workflowService) CreateWorkflow(username string, req requests.CreateWorkflowRequest) (*string, error) {
	workflowModel := models.Workflow{
		Name:  req.Name,
		Tasks: []models.Task{},
		Owner: username,
	}

	insertedID, err := service.workflowEntity.CreateWorkflow(workflowModel)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return insertedID, nil
}

func (service *workflowService) EditWorkflowByID(workflowID string, req requests.EditWorkflowRequest) (*models.Workflow, error) {
	workflowModel := models.Workflow{
		Name: req.Name,
	}

	workflow, err := service.workflowEntity.UpdateWorkflow(workflowID, workflowModel)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return workflow, nil
}

func (service *workflowService) DeleteWorkflowByID(workflowID string) error {
	if err := service.workflowEntity.DeleteWorkflow(workflowID); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (service *workflowService) TransferWorkflowByID(workflowID string, username string) (*models.Workflow, error) {
	workflowModel := models.Workflow{
		Owner: username,
	}

	workflow, err := service.workflowEntity.TransferWorkflowByID(workflowID, workflowModel)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return workflow, nil
}

func (service *workflowService) GetTasksByWorkflowID(workflowID string) ([]models.Task, error) {
	tasks, err := service.workflowEntity.FindTasksByWorkflowID(workflowID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return tasks, nil
}

func (service *workflowService) GetTaskByID(workflowID string, taskID string) (*models.Task, error) {
	task, err := service.workflowEntity.FindTaskByID(workflowID, taskID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return task, nil
}

func (service *workflowService) CreateTaskByWorkflowID(workflowID string, req requests.CreateTaskRequest) (*string, error) {
	var insertedID *string
	err := common.WithTransaction(context.TODO(), service.mongoClient, func(c context.Context, session mongo.Session) error {
		maxOrder, err := service.workflowEntity.FindMaxTaskOrderByWorkflowID(workflowID)
		if err != nil {
			return err
		}

		order := 1
		if maxOrder != nil {
			order = *maxOrder + 1
		}

		taskModel := models.Task{
			Name:        req.Name,
			Description: req.Description,
			Status:      models.TaskStatus("Pending"),
			Order:       order,
		}

		insertedID, err = service.workflowEntity.CreateTaskByWorkflowID(workflowID, taskModel)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return insertedID, nil
}

func (service *workflowService) EditTaskByID(workflowID string, taskID string, req requests.EditTaskRequest) (*models.Task, error) {
	taskModel := models.Task{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Order:       req.Order,
	}

	task, err := service.workflowEntity.UpdateTaskByID(workflowID, taskID, taskModel)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return task, nil
}

func (service *workflowService) DeleteTaskByID(workflowID string, taskID string) error {
	if err := service.workflowEntity.DeleteTaskByID(workflowID, taskID); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
