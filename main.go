//go:generate statik -src=frontend/build
//go:generate go fmt statik/statik.go
package main

import (
	"fmt"
	"github.com/vvisionnn/Drive-API/routers"
	"github.com/vvisionnn/Drive-API/settings"
	"log"
)

func main() {
	engine, err := routers.InitialRouter()
	if err != nil {
		log.Fatal(err)
	}

	port := settings.CONF.Port
	if err := engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("run server error: ", err)
	}
}
