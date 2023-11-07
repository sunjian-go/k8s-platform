package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
	"strconv"
)

var Workflow workflow

type workflow struct {
}

// 获取列表分页查询
func (w *workflow) GetWorkflows(c *gin.Context) {
	workflows := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Page      int    `form:"page"`
		Limit     int    `form:"limit"`
	})
	if err := c.Bind(&workflows); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Workflow.GetList(workflows.Name, workflows.Namespace, workflows.Page, workflows.Limit)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err,
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取wortflow列表成功",
		"data": data,
	})
}

// 查询workflow单条数据
func (w *workflow) GetById(c *gin.Context) {
	id := c.Query("id")
	idd, _ := strconv.Atoi(id)
	workflow, err := service.Workflow.GetById(idd)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err,
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "查询数据成功",
		"data": workflow,
	})
}

// 创建workflow
func (w *workflow) CreateWorkflow(c *gin.Context) {
	workcreate := new(service.WorkflowCreate)
	if err := c.ShouldBind(workcreate); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败: " + err.Error(),
		})
		return
	}

	err := service.Workflow.CreateWorkflow(workcreate)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "创建workflow失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "创建workflow成功",
	})
}

// 删除workflow
func (w *workflow) DelById(c *gin.Context) {
	id := c.Query("id")
	idd, _ := strconv.Atoi(id)
	fmt.Println("id=", id, "   idd=", idd)
	if err := service.Workflow.DelById(idd); err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除数据成功",
	})
}
