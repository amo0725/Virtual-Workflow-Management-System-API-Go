package responses

import (
	"virtual_workflow_management_system_gin/common"

	"github.com/gin-gonic/gin"
)

func Ok(ctx *gin.Context) {
	common.ResultJson(ctx, common.SUCCESS, "OK", map[string]interface{}{})
}

func OkWithData(ctx *gin.Context, data interface{}) {
	common.ResultJson(ctx, common.SUCCESS, "OK", data)
}
