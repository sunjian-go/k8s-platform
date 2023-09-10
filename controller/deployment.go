package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
	"strconv"
)

var Deployment deployment

type deployment struct {
}

// 获取deployment列表
func (d *deployment) GetDeployments(c *gin.Context) {
	deploy := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Limit     int    `form:"limit"`
		Page      int    `form:"page"`
	})
	if err := c.Bind(deploy); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	deployments, err := service.Deployment.GetDeployments(deploy.Name, deploy.Namespace, deploy.Limit, deploy.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  "获取deployment列表失败",
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取deployment列表成功",
		"data": deployments,
	})
}

// 获取deployment详情
func (d *deployment) GetDeploymentDetail(c *gin.Context) {
	deploy := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(deploy); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	fmt.Println("客户端传过来数据为：", deploy)
	deployment, err := service.Deployment.GetdeploymentDetail(deploy.Name, deploy.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  "获取deployment详情失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取deployment详情成功",
		"data": deployment,
	})
}

// 删除deployment
func (d *deployment) DeleteDeployment(c *gin.Context) {
	deploy := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(deploy); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	if err := service.Deployment.DeleteDeployment(deploy.Name, deploy.Namespace); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除deployment: " + deploy.Name + "成功",
	})
}

// 设置deployment副本数
func (d *deployment) ScaleDeployment(c *gin.Context) {
	deploy := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
		Replca    int    `form:"replca" `
	})

	if err := c.Bind(deploy); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	//fmt.Println("获取的数据为：", deploy)
	replca, err := service.Deployment.ScaleDeployment(deploy.Name, deploy.Namespace, deploy.Replca)
	if err != nil {
		c.JSON(400, gin.H{
			"err": "设置deployment副本失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "设置deployment: " + deploy.Name + "副本成功",
		"data": "当前副本数为：" + strconv.Itoa(int(replca)),
	})
}

// 创建deployment
func (d *deployment) CreateDeployment(c *gin.Context) {
	deployCreate := new(service.DeployCreate)
	if err := c.ShouldBind(deployCreate); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	fmt.Println("前端数据为：", deployCreate)

	err := service.Deployment.CreateDeployment(deployCreate)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "创建deployment成功",
	})
}

// 重启deployment
func (d *deployment) RestartDeployment(c *gin.Context) {
	deploy := new(struct {
		Name      string `json:"name" binding:"required"`
		Namespace string `json:"namespace" binding:"required"`
	})
	//PUT请求
	if err := c.ShouldBind(deploy); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	err := service.Deployment.RestartDeployment(deploy.Name, deploy.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "重启deployment成功",
	})
}

// 更新deployment
func (d *deployment) UpdateDeployment(c *gin.Context) {
	deploy := new(struct {
		Namespace string `json:"namespace" binding:"required"`
		Content   string `json:"content" binding:"required"`
	})
	//PUT请求
	if err := c.ShouldBindJSON(deploy); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	err := service.Deployment.UpdateDeployment(deploy.Namespace, deploy.Content)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新deployment成功",
	})
}

// 获取每个namespace中的deployment数量
func (d *deployment) GetNamespaceDeployNum(c *gin.Context) {
	deploys, err := service.Deployment.GetDeploymentNumPerNp()
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取deployment数量成功",
		"data": deploys,
	})
}
