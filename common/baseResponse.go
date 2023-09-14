package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	SUCCESS       = 0
	ERROR         = -1
	TOKEN_EXPIRED = -2
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Time    string      `json:"time"`
}

func ResultJson(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
		Time:    time.Now().Format(YYYYMMDDHHIISS),
	})
}
