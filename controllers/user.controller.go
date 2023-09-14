package controllers

import (
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/models"
	"virtual_workflow_management_system_gin/requests"
	"virtual_workflow_management_system_gin/responses"
	"virtual_workflow_management_system_gin/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.IUserService
}

func NewUserController(resource *databases.Resource) *UserController {
	userService := services.NewUserService(resource)
	return &UserController{UserService: userService}
}

// @Summary Register a new user
// @Tags Users
// @version 1.0
// @Description Register a new user with the input payload
// @Accept  application/json
// @Produce  application/json
// @Param user body requests.RegisterRequest true "User for registration"
// @Success 200 {object} string "OK"
// @Failure 400 {object} string "Invalid input"
// @Router /register [post]
func (controller *UserController) Register(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, "Invalid input")
		return
	}

	controller.UserService.Register(c, req)
}

// @Summary Login
// @Tags Users
// @version 1.0
// @Description Login with the input payload
// @Accept  application/json
// @Produce  application/json
// @Param user body requests.LoginRequest true "User for login"
// @Success 200 {object} string "OK"
// @Router /login [post]
func (controller *UserController) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, "Invalid input")
		return
	}

	controller.UserService.Login(c, req)
}

// @Summary Refresh token
// @Tags Users
// @version 1.0
// @Description Refresh token with the input payload
// @Accept  application/json
// @Produce  application/json
// @Param user body requests.RefreshTokenRequest true "User for refresh token"
// @Success 200 {object} string "OK"
// @Router /refresh-token [post]
func (controller *UserController) RefreshToken(c *gin.Context) {
	var req requests.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, "Invalid input")
		return
	}

	controller.UserService.RefreshToken(c, req)
}

// @Security access_token
// @Summary Logout
// @Tags Users
// @version 1.0
// @Description Logout
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} string "OK"
// @Router /logout [post]
func (controller *UserController) Logout(c *gin.Context) {
	user := c.MustGet("user").(models.JWTUser)

	err := controller.UserService.Logout(user.Username)
	if err != nil {
		responses.Error(c, err.Error())
		return
	}

	responses.Ok(c)
}
