package routers

import "github.com/gin-gonic/gin"
import "Drive-API/routers/monitor"

func InitialRouter() *gin.Engine {
	engine := gin.Default()

	api := engine.Group("api")
	{
		api.GET("/ping", monitor.Ping)
	}

	return engine
}