package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var Ingress ingress

type ingress struct {
}

// 获取ingress列表
func (i *ingress) GetIngresses(c *gin.Context) {
	ing := new(struct {
		Name      string `form:"name" `
		Namespace string `form:"namespace" `
		Limit     int    `form:"limit" `
		Page      int    `form:"page" `
	})
	if err := c.Bind(ing); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	ingresses, err := service.Ingress.GetIngresses(ing.Name, ing.Namespace, ing.Limit, ing.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取ingress列表成功",
		"data": ingresses,
	})
}

// 获取Ingress详情
func (i *ingress) GetIngressDetail(c *gin.Context) {
	ing := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(ing); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	ingress, err := service.Ingress.GetIngDetail(ing.Name, ing.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取ingress详情成功",
		"data": ingress,
	})
}

// 创建Ingress
func (i *ingress) CreateIngress(c *gin.Context) {
	ing := service.IngressCreate{}
	if err := c.ShouldBind(&ing); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败: " + err.Error(),
		})
		return
	}
	if err := service.Ingress.CreateIngress(&ing); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "创建ingress成功",
	})
}

// 删除Ingress
func (i *ingress) DeleteIngress(c *gin.Context) {
	ing := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(ing); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败: " + err.Error(),
		})
		return
	}
	if err := service.Ingress.DeleteIngress(ing.Name, ing.Namespace); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(20, gin.H{
		"msg": "删除ingress成功",
	})

}

// 更新Ingress
func (i *ingress) UpdateIngress(c *gin.Context) {
	ing := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	if err := c.ShouldBindJSON(ing); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败: " + err.Error(),
		})
		return
	}
	if err := service.Ingress.UpdateIngress(ing.Namespace, ing.Content); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新ingress成功",
	})
}
