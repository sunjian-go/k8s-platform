package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var Node node

type node struct {
}

func (n *node) GetNodes(c *gin.Context) {
	//GET请求
	nodes, err := service.Node.GetNodes()
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取node列表成功",
		"data": nodes,
	})
}

// 获取node详情
func (n *node) GetNodeDetail(c *gin.Context) {
	//POST请求
	node := new(struct {
		Name string `form:"name"`
	})
	//nodeName := c.Params.ByName("name")
	if err := c.Bind(node); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}

	fmt.Println("nodeName为： ", node.Name)
	nodedetail, err := service.Node.GetNodeDetail(node.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取node详情成功",
		"data": nodedetail,
	})
}
