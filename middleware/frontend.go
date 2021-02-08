package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/pkgs/frontend"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// FrontendHandler handle front static file, largely inspired by Cloudreve, thanks HFO4
// https://github.com/cloudreve/Cloudreve
func FrontendHandler() gin.HandlerFunc {
	indexFile, err := frontend.StaticFS.Open("/index.html")
	if err != nil {
		log.Panic(err)
	}
	fileContent, err := ioutil.ReadAll(indexFile)
	if err != nil {
		log.Panic(err)
	}
	fileContentStr := string(fileContent)

	fileServer := http.FileServer(frontend.StaticFS)
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// API
		if strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}

		if (path == "/index.html") || (path == "/") || !frontend.StaticFS.Exists("/", path) {
			c.Header("Content-Type", "text/html")
			c.String(200, fileContentStr)
			c.Abort()
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}
