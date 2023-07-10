package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var Namespace namespace

type namespace struct {
}

// 获取namespace列表
func (n *namespace) GetNamespaces(c *gin.Context) {
	namespaces, err := service.Namespace.GetNamespaces()
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取namespace列表成功",
		"data": namespaces,
	})
}

// 获取Namespace详情
func (n *namespace) GetNamespaceDetail(c *gin.Context) {
	//GET请求
	namespace := new(struct {
		Name string `form:"name"`
	})
	if err := c.Bind(namespace); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	namespaceDetail, err := service.Namespace.GetNamespaceDetail(namespace.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取namespace详情成功",
		"data": namespaceDetail,
	})
}

// 删除Namespace
func (n *namespace) DeleteNamespace(c *gin.Context) {
	//Delete请求
	//namespaceName := c.Param("name")
	namespace := new(struct {
		Name string `form:"name"`
	})
	if err := c.Bind(namespace); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	//fmt.Println("namespace: ", namespace.Name)
	if err := service.Namespace.DeleteNamespace(namespace.Name); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除namespace " + namespace.Name + " 成功",
	})
}
