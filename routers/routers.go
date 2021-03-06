package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/middleware"
	"github.com/vvisionnn/Drive-API/routers/drive"
	"github.com/vvisionnn/Drive-API/routers/monitor"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
)

func InitialRouter() (*gin.Engine, error) {
	engine := gin.Default()
	store := persistence.NewInMemoryStore(time.Second)
	if err := drive.InitialDrive(); err != nil {
		return nil, err
	}

	engine.Use(middleware.FrontendHandler())

	api := engine.Group("api")
	{
		api.GET("/ping", monitor.Ping)
		api.GET("/cache_ping", cache.CachePage(store, time.Second*10, monitor.CachePing))

		auth := api.Group("auth")
		{
			auth.PUT("/conf", drive.SetConfiguration)
			auth.GET("/stat", drive.StatusHandler)
			//auth.GET("/url", drive.UrlHandler)
			auth.GET("/callback", drive.CallbackHandler)
		}

		driveApi := api.Group("drive")
		{
			driveApi.GET("", cache.CachePage(store, time.Second*5, drive.ListRootHandler))
			driveApi.GET("/:id", cache.CachePage(store, time.Second*10, drive.ListHandler))
		}

		betaApi := api.Group("beta")
		{
			betaApi.GET("/cache", monitor.CacheHandler)
		}
	}

	return engine, nil
}
