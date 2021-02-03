package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type respContent struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InternalServerError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, respContent{
		Code:    0,
		Message: msg,
		Data:    nil,
	})
}

func Success(ctx *gin.Context) {
	SuccessWithMessageAndData(ctx, "success", nil)
}

func SuccessWithData(ctx *gin.Context, data interface{}) {
	SuccessWithMessageAndData(ctx, "success", data)
}

func SuccessWithMessage(ctx *gin.Context, msg string) {
	SuccessWithMessageAndData(ctx, msg, nil)
}

func SuccessWithMessageAndData(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, respContent{
		Code:    1,
		Message: msg,
		Data:    data,
	})
}


func RedirectTemporary(ctx *gin.Context, url string) {
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

