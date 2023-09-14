package repositories

import (
	"context"
	"errors"
	"virtual_workflow_management_system_gin/common"
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var WorkflowEntity IWorkflow

type workflowEntity struct {
	resource    *databases.Resource
	repository  *mongo.Collection
	mongoClient *mongo.Client
}

type IWorkflow interface {
	FindWorkflowsByUsername(username string) ([]models.Workflow, error)
	FindWorkflowByID(workflowID string) (*models.Workflow, error)
	CreateWorkflow(workflow models.Workflow) (*string, error)
	UpdateWorkflow(workflowID string, workflow models.Workflow) (*models.Workflow, error)
	DeleteWorkflow(workflowID string) error
	TransferWorkflowByID(workflowID string, workflow models.Workflow) (*models.Workflow, error)
	FindTasksByWorkflowID(workflowID string) ([]models.Task, error)
	FindTaskByID(workflowID string, taskID string) (*models.Task, error)
	FindMaxTaskOrderByWorkflowID(workflowID string) (*int, error)
	CreateTaskByWorkflowID(workflowID string, task models.Task) (*string, error)
	UpdateTaskByID(workflowID string, taskID string, task models.Task) (*models.Task, error)
	DeleteTaskByID(workflowID string, taskID string) error
}

func NewWorkflowEntity(resource *databases.Resource) IWorkflow {
	workflowRepository := resource.MongoDB.Collection("workflows")
	WorkflowEntity = &workflowEntity{resource: resource, repository: workflowRepository, mongoClient: resource.MongoDB.Client()}
	return WorkflowEntity
}

func (entity *workflowEntity) FindWorkflowsByUsername(username string) ([]models.Workflow, error) {
	ctx, cancel := initContext()
	defer cancel()

	cursor, err := entity.repository.Find(ctx, bson.M{
		"owner": username,
	})
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("failed to retrieve workflows")
	}

	var workflows []models.Workflow
	err = cursor.All(ctx, &workflows)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("failed to retrieve workflows")
	}

	return workflows, nil
}

func (entity *workflowEntity) FindWorkflowByID(workflowID string) (*models.Workflow, error) {
	ctx, cancel := initContext()
	defer cancel()

	workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("invalid ObjectID format")
	}

	filter := bson.M{"_id": workflowObjectID}
	var workflow models.Workflow
	err = entity.repository.FindOne(ctx, filter).Decode(&workflow)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("workflow does not exist")
	}

	return &workflow, nil
}

func (entity *workflowEntity) CreateWorkflow(workflow models.Workflow) (*string, error) {
	ctx, cancel := initContext()
	defer cancel()

	workflow.SetCreatedAt()
	workflow.SetUpdatedAt()

	insertResult, err := entity.repository.InsertOne(ctx, workflow)
	if err != nil {
		logrus.Errorf("Failed to insert new workflow: %v", err)
		return nil, errors.New("failed to create workflow")
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		logrus.Error("Failed to convert InsertedID to ObjectID")
		return nil, errors.New("failed to convert InsertedID to ObjectID")
	}

	insertedIDString := insertedID.Hex()

	return &insertedIDString, nil
}

func (entity *workflowEntity) UpdateWorkflow(workflowID string, workflow models.Workflow) (*models.Workflow, error) {
	var updatedWorkflow models.Workflow
	err := common.WithTransaction(context.TODO(), entity.mongoClient, func(c context.Context, session mongo.Session) error {
		ctx, cancel := initContext()
		defer cancel()

		workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		workflow.SetUpdatedAt()

		filter := bson.M{"_id": workflowObjectID}
		update := bson.M{
			"$set": bson.M{
				"name":       workflow.Name,
				"updated_at": workflow.UpdatedAt,
			},
		}

		result, err := entity.repository.UpdateOne(ctx, filter, update)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to update workflow")
		}

		if result.ModifiedCount == 0 {
			return errors.New("no workflow was updated")
		}

		err = entity.repository.FindOne(ctx, filter).Decode(&updatedWorkflow)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to retrieve updated workflow")
		}

		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &updatedWorkflow, nil
}

func (entity *workflowEntity) DeleteWorkflow(workflowID string) error {
	err := common.WithTransaction(context.TODO(), entity.mongoClient, func(c context.Context, session mongo.Session) error {
		ctx, cancel := initContext()
		defer cancel()

		workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		filter := bson.M{"_id": workflowObjectID}

		_, err = entity.repository.DeleteOne(ctx, filter)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to delete workflow")
		}

		return nil
	})
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (entity *workflowEntity) TransferWorkflowByID(workflowID string, workflow models.Workflow) (*models.Workflow, error) {
	var updatedWorkflow models.Workflow
	err := common.WithTransaction(context.TODO(), entity.mongoClient, func(c context.Context, session mongo.Session) error {
		ctx, cancel := initContext()
		defer cancel()

		workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		workflow.SetUpdatedAt()

		filter := bson.M{"_id": workflowObjectID}
		update := bson.M{
			"$set": bson.M{
				"owner": workflow.Owner,
			},
		}

		result, err := entity.repository.UpdateOne(ctx, filter, update)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to update workflow")
		}

		if result.ModifiedCount == 0 {
			return errors.New("no workflow was updated")
		}

		err = entity.repository.FindOne(ctx, filter).Decode(&updatedWorkflow)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to retrieve updated workflow")
		}

		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &updatedWorkflow, nil
}

func (entity *workflowEntity) FindTasksByWorkflowID(workflowID string) ([]models.Task, error) {
	ctx, cancel := initContext()
	defer cancel()

	workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("invalid ObjectID format")
	}

	aggregatePipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"_id": workflowObjectID}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$tasks", "preserveNullAndEmptyArrays": false}}},
		bson.D{{Key: "$sort", Value: bson.M{"tasks.order": 1}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$_id", "tasks": bson.M{"$push": "$tasks"}}}},
	}

	cursor, err := entity.repository.Aggregate(ctx, aggregatePipeline)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("failed to aggregate tasks")
	}

	var results []struct {
		ID    primitive.ObjectID `bson:"_id"`
		Tasks []models.Task      `bson:"tasks"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		logrus.Error(err)
		return nil, errors.New("failed to decode tasks")
	}

	if len(results) == 0 {
		return nil, errors.New("no tasks found")
	}

	return results[0].Tasks, nil
}

func (entity *workflowEntity) FindTaskByID(workflowID string, taskID string) (*models.Task, error) {
	_, cancel := initContext()
	defer cancel()

	taskObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("invalid ObjectID format")
	}

	workflow, err := WorkflowEntity.FindWorkflowByID(workflowID)

	if err != nil {
		logrus.Error(err)
		return nil, errors.New("workflow does not exist")
	}

	for _, task := range workflow.Tasks {
		if task.ID == taskObjectID {
			return &task, nil
		}
	}

	return nil, errors.New("task does not exist")
}

func (entity *workflowEntity) FindMaxTaskOrderByWorkflowID(workflowID string) (*int, error) {
	ctx, cancel := initContext()
	defer cancel()

	workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("invalid ObjectID format")
	}

	aggregatePipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"_id": workflowObjectID}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$tasks", "preserveNullAndEmptyArrays": false}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": "$_id", "maxOrder": bson.M{"$max": "$tasks.order"}}}},
	}

	cursor, err := entity.repository.Aggregate(ctx, aggregatePipeline)
	if err != nil {
		logrus.Error("Failed to aggregate: ", err)
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		logrus.Error("Failed to decode cursor: ", err)
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	maxOrder := results[0]["maxOrder"]
	maxOrderInt := int(maxOrder.(int32))

	return &maxOrderInt, nil
}

func (entity *workflowEntity) CreateTaskByWorkflowID(workflowID string, task models.Task) (*string, error) {
	var taskIDString string
	err := common.WithTransaction(context.TODO(), entity.mongoClient, func(c context.Context, session mongo.Session) error {
		ctx, cancel := initContext()
		defer cancel()

		workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		task.ID = primitive.NewObjectID()
		task.SetCreatedAt()
		task.SetUpdatedAt()

		filter := bson.M{"_id": workflowObjectID}
		update := bson.M{
			"$push": bson.M{
				"tasks": task,
			},
		}

		updateResult, err := entity.repository.UpdateOne(ctx, filter, update)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to create task")
		}

		if updateResult.MatchedCount == 0 {
			logrus.Error("workflow not found")
			return errors.New("workflow not found")
		}

		var updatedWorkflow models.Workflow
		err = entity.repository.FindOne(ctx, filter).Decode(&updatedWorkflow)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to retrieve updated workflow")
		}

		updatedTasks := updatedWorkflow.Tasks
		updatedTasksLength := len(updatedTasks)

		if updatedTasksLength == 0 {
			return errors.New("no task was created")
		}

		taskIDString = updatedTasks[updatedTasksLength-1].ID.Hex()

		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &taskIDString, nil
}

func (entity *workflowEntity) UpdateTaskByID(workflowID string, taskID string, task models.Task) (*models.Task, error) {
	var updatedTaskModel models.Task
	err := common.WithTransaction(context.TODO(), entity.mongoClient, func(c context.Context, session mongo.Session) error {
		ctx, cancel := initContext()
		defer cancel()

		workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		taskObjectID, err := primitive.ObjectIDFromHex(taskID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		task.SetUpdatedAt()

		filter := bson.M{"_id": workflowObjectID, "tasks._id": taskObjectID}
		update := bson.M{
			"$set": bson.M{
				"tasks.$.name":        task.Name,
				"tasks.$.description": task.Description,
				"tasks.$.status":      task.Status,
				"tasks.$.order":       task.Order,
				"tasks.$.updated_at":  task.UpdatedAt,
			},
		}

		updateResult, err := entity.repository.UpdateOne(ctx, filter, update)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to update task")
		}

		if updateResult.MatchedCount == 0 {
			return errors.New("no task was updated")
		}

		var updatedWorkflow models.Workflow
		err = entity.repository.FindOne(ctx, filter).Decode(&updatedWorkflow)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to retrieve updated workflow")
		}

		updatedTasks := updatedWorkflow.Tasks
		for _, updatedTask := range updatedTasks {
			if updatedTask.ID == taskObjectID {
				updatedTaskModel = updatedTask
				return nil
			}
		}

		return errors.New("task does not exist")
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &updatedTaskModel, nil
}

func (entity *workflowEntity) DeleteTaskByID(workflowID string, taskID string) error {
	err := common.WithTransaction(context.TODO(), entity.mongoClient, func(c context.Context, session mongo.Session) error {
		ctx, cancel := initContext()
		defer cancel()

		workflowObjectID, err := primitive.ObjectIDFromHex(workflowID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		taskObjectID, err := primitive.ObjectIDFromHex(taskID)
		if err != nil {
			logrus.Error(err)
			return errors.New("invalid ObjectID format")
		}

		filter := bson.M{"_id": workflowObjectID}
		update := bson.M{
			"$pull": bson.M{
				"tasks": bson.M{
					"_id": taskObjectID,
				},
			},
		}

		updateResult, err := entity.repository.UpdateOne(ctx, filter, update)
		if err != nil {
			logrus.Error(err)
			return errors.New("failed to delete task")
		}

		if updateResult.MatchedCount == 0 {
			return errors.New("no task was deleted")
		}

		return nil
	})
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
