package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var StatefulSet statefulSet

type statefulSet struct {
}

// 获取statefulSet列表
func (s *statefulSet) GetStatefulSets(c *gin.Context) {
	//GET请求
	statefulset := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
		Limit     int    `form:"limit" binding:"required"`
		Page      int    `form:"page" binding:"required"`
	})
	if err := c.Bind(statefulset); err != nil {
		c.JSON(400, gin.H{
			"err":  "数据绑定失败",
			"data": nil,
		})
		return
	}
	statefulsets, err := service.StatefulSet.GetStatefulSets(statefulset.Name, statefulset.Namespace, statefulset.Limit, statefulset.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取statefulSet列表成功",
		"data": statefulsets,
	})
}

// 获取statefulSet详情
func (s *statefulSet) GetStatefulSetDetail(c *gin.Context) {
	//GET请求
	statefulset := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(statefulset); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	statefulsetdetail, err := service.StatefulSet.GetStatefulSetDetail(statefulset.Name, statefulset.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取statefulSet详情成功",
		"data": statefulsetdetail,
	})
}

// 删除statefulSet
func (s *statefulSet) DeleteStatefulSet(c *gin.Context) {
	//DELETE请求
	statefulset := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(statefulset); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	err := service.StatefulSet.DeleteStatefulSet(statefulset.Name, statefulset.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除statefulSet成功",
	})
}

// 更新statefulSet
func (s *statefulSet) UpdateStatefulSet(c *gin.Context) {
	//PUT请求
	statefulset := new(struct {
		Namespace string `json:"namespace" binding:"required"`
		Content   string `json:"content" binding:"required"`
	})
	if err := c.ShouldBindJSON(statefulset); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	err := service.StatefulSet.UpdateStatefulSet(statefulset.Namespace, statefulset.Content)
	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新statefulSet成功",
	})
}
