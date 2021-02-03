package main

import (
	"Drive-API/routers"
	"Drive-API/settings"
	"fmt"
	"log"
)

func main() {
	engine := routers.InitialRouter()

	port := settings.CONF.Port
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("run server error: ", err)
	}
}
