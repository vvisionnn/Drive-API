package routers

import (
	"Drive-API/routers/drive"
	"github.com/gin-gonic/gin"
)
import "Drive-API/routers/monitor"

func InitialRouter() *gin.Engine {
	engine := gin.Default()

	api := engine.Group("api")
	{
		api.GET("/ping", monitor.Ping)

		auth := api.Group("auth")
		{
			auth.GET("/stat", drive.StatusHandler)
			auth.GET("/url", drive.UrlHandler)
			auth.GET("/callback", drive.CallbackHandler)
		}

		api.GET("/drive", drive.ListHandler)
	}

	return engine
}