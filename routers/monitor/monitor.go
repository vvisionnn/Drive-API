package monitor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/pkgs/onedrive"
	"github.com/vvisionnn/Drive-API/pkgs/response"
	"github.com/vvisionnn/Drive-API/routers/drive"
	"sync"
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
	wg := sync.WaitGroup{}
	for _, item := range items.Value {
		wg.Add(1)
		go func(i onedrive.ItemInfo) {
			node := &Node{
				i,
				[]*Node{},
			}
			node.BFSGetChildren(2)
			result = append(result, node)
			wg.Done()
		}(item)
	}
	wg.Wait()
	response.Success(ctx)
}

type Node struct {
	onedrive.ItemInfo
	children []*Node
}

func (node *Node) BFSGetChildren(depth int) {
	fmt.Println("depth: ", depth, "id: ", node.ID, "name: ", node.Name)
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
	if depth--; depth == 0 || len(node.children) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	for _, child := range node.children {
		wg.Add(1)
		go func(c *Node) {
			c.BFSGetChildren(depth)
			wg.Done()
		}(child)
	}
	wg.Wait()
}
