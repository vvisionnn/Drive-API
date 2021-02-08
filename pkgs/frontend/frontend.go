package frontend

import (
	"github.com/gin-contrib/static"
	"github.com/rakyll/statik/fs"
	_ "github.com/vvisionnn/Drive-API/statik"
	"log"
	"net/http"
)

type GinFS struct {
	FS http.FileSystem
}

var StaticFS static.ServeFileSystem

// Open 打开文件
func (b *GinFS) Open(name string) (http.File, error) {
	return b.FS.Open(name)
}

// Exists 文件是否存在
func (b *GinFS) Exists(prefix string, filepath string) bool {
	if _, err := b.FS.Open(filepath); err != nil {
		return false
	}
	return true
}

func init() {
	var err error
	StaticFS = &GinFS{}
	StaticFS.(*GinFS).FS, err = fs.New()
	if err != nil {
		log.Panic(err)
	}
}
