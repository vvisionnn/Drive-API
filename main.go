//go:generate statik -src=frontend/build
//go:generate go fmt statik/statik.go
package main

import (
	"github.com/vvisionnn/Drive-API/routers"
	"log"
)

func main() {
	engine, err := routers.InitialRouter()
	if err != nil {
		log.Fatal(err)
	}

	if err := engine.Run(":8421"); err != nil {
		log.Fatal("run server error: ", err)
	}
}
