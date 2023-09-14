package services

import (
	"errors"
	"virtual_workflow_management_system_gin/common"
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/middlewares"
	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/repositories"
	"virtual_workflow_management_system_gin/requests"
	"virtual_workflow_management_system_gin/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	userEntity repositories.IUser
	redis      *redis.Client
}

func NewUserService(resource *databases.Resource) *UserService {
	return &UserService{
		userEntity: repositories.NewUserEntity(resource),
		redis:      resource.Redis,
	}
}

func (service *UserService) Register(c *gin.Context, req requests.RegisterRequest) {
	reqWithHashedPassword := requests.RegisterRequest{
		Username: req.Username,
		Password: common.HashPassword(req.Password),
		Role:     req.Role,
	}

	_, err := service.userEntity.CreateOne(reqWithHashedPassword)
	if err != nil {
		logrus.Error(err)
		responses.Error(c, err.Error())
		return
	}

	responses.Ok(c)
}

func (service *UserService) Login(c *gin.Context, req requests.LoginRequest) {
	user, err := service.userEntity.FindOneByUsername(req.Username)
	if err != nil {
		logrus.Error(err)
		responses.Error(c, err.Error())
		return
	}

	if user == nil {
		responses.Error(c, "username does not exist")
		return
	}

	if common.ComparePasswordAndHashedPassword(req.Password, user.Password) != nil {
		responses.Error(c, "wrong password")
		return
	}

	jwt, err := middlewares.GenerateJWTToken(*user, service.redis)
	if err != nil {
		logrus.Error(err)
		responses.Error(c, "failed to generate token")
		return
	}

	responses.OkWithData(c, gin.H{
		"access_token":  jwt["access_token"],
		"refresh_token": jwt["refresh_token"],
	})
}

func (service *UserService) RefreshToken(c *gin.Context, req requests.RefreshTokenRequest) {
	jwt, err := middlewares.RefreshJWTToken(req.RefreshToken, service.redis)
	if err != nil {
		logrus.Error(err)
		responses.Error(c, "failed to refresh token")
		return
	}

	responses.OkWithData(c, gin.H{
		"access_token":  jwt["access_token"],
		"refresh_token": jwt["refresh_token"],
	})
}

func (service *UserService) Logout(username string) error {

	err := middlewares.DeleteJWTToken(username, service.redis)
	if err != nil {
		logrus.Error(err)
		return errors.New("failed to logout")
	}

	return nil
}

func (service *UserService) GetUsersByUsername(username string) (*models.User, error) {
	user, err := service.userEntity.FindOneByUsername(username)
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("failed to get user")
	}

	return user, nil
}
