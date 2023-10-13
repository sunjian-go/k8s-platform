package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var SVC svc

type svc struct {
}

// 获取svc列表
func (s *svc) GetSvcs(c *gin.Context) {
	//GET请求
	svc := new(struct {
		Name      string `form:"name" `
		Namespace string `form:"namespace" `
		Limit     int    `form:"limit" `
		Page      int    `form:"page" `
	})
	if err := c.Bind(svc); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	svcresp, err := service.SVC.GetSvcs(svc.Name, svc.Namespace, svc.Limit, svc.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取svc列表成功",
		"data": svcresp,
	})
}

// 获取Service详情
func (s *svc) GetSvcDetail(c *gin.Context) {
	//GET请求
	svc := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(svc); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	svcdetail, err := service.SVC.GetSvcDetail(svc.Name, svc.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取svc详情成功",
		"data": svcdetail,
	})
}

// 创建Service
func (s *svc) CreateSvc(c *gin.Context) {
	//POST请求
	svc := new(service.ServiceCreate)
	if err := c.ShouldBind(svc); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	if err := service.SVC.CreateSvc(svc); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "创建svc成功",
	})
}

// 删除Service
func (s *svc) DeleteSvc(c *gin.Context) {
	//DELETE请求
	svc := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(svc); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败: " + err.Error(),
			"data": nil,
		})
		return
	}
	err := service.SVC.DeleteSvc(svc.Name, svc.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取svc详情成功",
		"data": svc,
	})
}

// 更新Service
func (s *svc) UpdateSvc(c *gin.Context) {
	//PUT请求
	svc := new(struct {
		Namespace string `json:"namespace" binding:"required"`
		Centent   string `json:"centent" binding:"required"`
	})
	if err := c.ShouldBindJSON(svc); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	err := service.SVC.UpdateSvc(svc.Namespace, svc.Centent)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除svc成功",
	})
}
