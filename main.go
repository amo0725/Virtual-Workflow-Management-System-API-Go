package main

import (
	"virtual_workflow_management_system_gin/databases"
	"virtual_workflow_management_system_gin/middlewares"
	"virtual_workflow_management_system_gin/routes"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @securityDefinitions.apikey access_token
// @in header
// @name Authorization
// @securitySchemes: access_token
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error(err)
	}
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(middlewares.NewCors([]string{"*"}))
	r.GET("swagger/*any", middlewares.NewSwagger())
	publicRoute := r.Group(os.Getenv("BASE_PATH"))
	resource, err := databases.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()
	routes.InitUserRouter(publicRoute, resource)
	routes.InitWorkflowRouter(publicRoute, resource)
	r.Run(":" + os.Getenv("PORT"))
}
