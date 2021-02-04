package monitor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/pkgs/response"
	"time"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithMessage(ctx, fmt.Sprintf("pong: %d", time.Now().Unix()))
}

func CachePing(ctx *gin.Context) {
	response.SuccessWithMessage(ctx, fmt.Sprintf("pong: %d", time.Now().Unix()))
}
