package responses

import (
	"virtual_workflow_management_system_gin/common"

	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, msg string) {
	common.ResultJson(ctx, common.ERROR, msg, map[string]interface{}{})
}

func ErrorWithToken(ctx *gin.Context, msg string) {
	common.ResultJson(ctx, common.TOKEN_EXPIRED, msg, map[string]interface{}{})
}
