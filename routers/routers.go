package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/routers/drive"
	"github.com/vvisionnn/Drive-API/routers/monitor"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
)

func InitialRouter() *gin.Engine {
	engine := gin.Default()
	store := persistence.NewInMemoryStore(time.Second)

	api := engine.Group("api")
	{
		api.GET("/ping", monitor.Ping)
		api.GET("/cache_ping", cache.CachePage(store, time.Second * 10, monitor.CachePing))

		auth := api.Group("auth")
		{
			auth.GET("/stat", drive.StatusHandler)
			auth.GET("/url", drive.UrlHandler)
			auth.GET("/callback", drive.CallbackHandler)
		}

		api.GET("/drive", cache.CachePage(store, time.Second * 5, drive.ListHandler))
	}

	return engine
}