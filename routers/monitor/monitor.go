package monitor

import (
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/pkgs/response"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithMessage(ctx, "pong")
}

