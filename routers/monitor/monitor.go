package monitor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/pkgs/onedrive"
	"github.com/vvisionnn/Drive-API/pkgs/response"
	"github.com/vvisionnn/Drive-API/routers/drive"
	"time"
)

func Ping(ctx *gin.Context) {
	response.SuccessWithMessage(ctx, fmt.Sprintf("pong: %d", time.Now().Unix()))
}

func CachePing(ctx *gin.Context) {
	response.SuccessWithMessage(ctx, fmt.Sprintf("pong: %d", time.Now().Unix()))
}

// api for test
func CacheHandler(ctx *gin.Context) {
	// recursive get folder
	// list root, get all folder's id
	// BFS get children
	items, err := drive.Drive.ListRootChildren()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	var result []*Node
	for _, item := range items.Value {
		node := &Node{
			item,
			[]*Node{},
		}
		node.BFSGetChildren()
		result = append(result, node)
	}
	fmt.Println(result)
	response.Success(ctx)
}

type Node struct {
	onedrive.ItemInfo
	children []*Node
}

func (node *Node) BFSGetChildren() {
	fmt.Println("id: ", node.ID, "name: ", node.Name)
	items, err := drive.Drive.ListItemChildren(node.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, item := range items.Value {
		node.children = append(node.children, &Node{
			item,
			[]*Node{},
		})
	}

	for _, child := range node.children {
		child.BFSGetChildren()
	}
}
