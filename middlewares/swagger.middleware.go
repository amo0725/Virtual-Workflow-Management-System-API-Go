package middlewares

import (
	"os"
	"virtual_workflow_management_system_gin/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewSwagger() gin.HandlerFunc {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error(err)
	}

	docs.SwaggerInfo.Title = "Virtual Workflow Management System API"
	docs.SwaggerInfo.Description = "This is a server for Virtual Workflow Management System."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = os.Getenv("HOST") + ":" + os.Getenv("PORT")
	docs.SwaggerInfo.BasePath = os.Getenv("BASE_PATH")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
