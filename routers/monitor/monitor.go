package monitor

import (
	"Drive-API/pkgs/response"
	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithMessage(ctx, "pong")
}

