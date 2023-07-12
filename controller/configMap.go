package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var ConfigMap configMap

type configMap struct {
}

// 获取cm列表
func (c *configMap) GetConfigMaps(ctx *gin.Context) {
	cm := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Limit     int    `form:"limit"`
		Page      int    `form:"page"`
	})
	if err := ctx.Bind(cm); err != nil {
		ctx.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	cms, err := service.ConfigMap.GetConfigMaps(cm.Name, cm.Namespace, cm.Limit, cm.Page)
	if err != nil {
		ctx.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"msg":  "获取cm列表成功",
		"data": cms,
	})
}

// 获取cm详情
func (c *configMap) GetConfigDetail(ctx *gin.Context) {
	cm := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := ctx.Bind(cm); err != nil {
		ctx.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	cmdetail, err := service.ConfigMap.GetConfigMapDetail(cm.Name, cm.Namespace)
	if err != nil {
		ctx.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"msg":  "获取cm详情成功",
		"data": cmdetail,
	})
}

// 删除cm
func (c *configMap) DeleteConfigMap(ctx *gin.Context) {
	cm := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := ctx.Bind(cm); err != nil {
		ctx.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	if err := service.ConfigMap.DeleteConfigMap(cm.Name, cm.Namespace); err != nil {
		ctx.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"msg": "删除cm成功",
	})
}

// 更新cm
func (c *configMap) UpdateConfigMap(ctx *gin.Context) {
	cm := new(struct {
		Namespace string `json:"namespace" binding:"required"`
		Content   string `json:"content" binding:"required"`
	})
	if err := ctx.ShouldBindJSON(cm); err != nil {
		ctx.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	if err := service.ConfigMap.UpdateConfigMap(cm.Namespace, cm.Content); err != nil {
		ctx.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"msg": "更新cm成功",
	})
}
