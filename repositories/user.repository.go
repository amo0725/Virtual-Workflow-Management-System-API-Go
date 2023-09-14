package repositories

import (
	"errors"
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/requests"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserEntity IUser

type userEntity struct {
	resource   *databases.Resource
	repository *mongo.Collection
}

type IUser interface {
	CreateOne(user requests.RegisterRequest) (*models.User, error)
	FindOneByUsername(username string) (*models.User, error)
}

func NewUserEntity(resource *databases.Resource) IUser {
	if resource == nil || resource.MongoDB == nil || resource.Redis == nil {
		return &userEntity{}
	}
	userRepository := resource.MongoDB.Collection("users")
	UserEntity = &userEntity{resource: resource, repository: userRepository}
	return UserEntity
}

func (entity *userEntity) CreateOne(user requests.RegisterRequest) (*models.User, error) {
	ctx, cancel := initContext()
	defer cancel()

	userModel := models.User{
		Username: user.Username,
		Password: user.Password,
	}

	existingUser, err := entity.FindOneByUsername(user.Username)

	if err != nil && err.Error() != "username does not exist" {
		logrus.Error(err)
		return nil, errors.New("internal Server Error")
	}

	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	_, err = entity.repository.InsertOne(ctx, userModel)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("failed to create user")
	}

	return &userModel, nil
}

func (entity *userEntity) FindOneByUsername(username string) (*models.User, error) {
	ctx, cancel := initContext()
	defer cancel()

	filter := bson.M{"username": username}
	var user models.User
	err := entity.repository.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("username does not exist")
	}

	return &user, nil
}
