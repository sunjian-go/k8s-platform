package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var DaemonSet daemonSet

type daemonSet struct {
}

// 获取daemonSet列表
func (d *daemonSet) GetDaemonSets(c *gin.Context) {
	daemonset := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Limit     int    `form:"limit" binding:"required"`
		Page      int    `form:"page" binding:"required"`
	})
	if err := c.Bind(daemonset); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	daemonsets, err := service.DaemonSet.GetDaemonSets(daemonset.Name, daemonset.Namespace, daemonset.Limit, daemonset.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取daemonSet列表成功",
		"data": daemonsets,
	})
}

// 获取daemonset详情
func (d *daemonSet) GetDaemonSetDetail(c *gin.Context) {
	daemonset := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(daemonset); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	daemonsetDetail, err := service.DaemonSet.GetDaemonSetDetail(daemonset.Name, daemonset.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取daemonSet详情成功",
		"data": daemonsetDetail,
	})
}

// 删除daemonset
func (d *daemonSet) DeleteDaemonSet(c *gin.Context) {
	daemonset := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	//Delete请求
	if err := c.Bind(daemonset); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.DaemonSet.DeleteDaemonSet(daemonset.Name, daemonset.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "删除daemonSet成功",
		"data": nil,
	})
}

// 更新daemonset
func (d *daemonSet) UpdateDaemonSet(c *gin.Context) {
	//PUT请求
	daemonset := new(struct {
		Namespace string `json:"namespace" binding:"required"`
		Content   string `json:"content" binding:"required"`
	})
	if err := c.ShouldBindJSON(daemonset); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败: " + err.Error(),
		})
		return
	}
	err := service.DaemonSet.UpdateDaemonSet(daemonset.Namespace, daemonset.Content)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新daemonSet成功",
	})
}

// 创建daemonset
func (d *daemonSet) CreateDaemonSet(c *gin.Context) {
	daemonset := new(service.DaemonSetCreate)
	if err := c.ShouldBind(daemonset); err != nil {
		c.JSON(400, gin.H{
			"err":  "数据绑定失败：" + err.Error(),
			"data": "",
		})
		return
	}
	err := service.DaemonSet.CreateDaemonSet(daemonset)
	if err != nil {
		c.JSON(400, gin.H{
			"msg":  err.Error(),
			"data": "",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "创建daemonset成功",
		"data": "",
	})
}
