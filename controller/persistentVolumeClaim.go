package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var Pvc pvc

type pvc struct {
}

// 获取pvc列表
func (p *pvc) GetPvcs(c *gin.Context) {
	pvc := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Limit     int    `form:"limit"`
		Page      int    `form:"page"`
	})
	if err := c.Bind(pvc); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	pvcresp, err := service.Pvc.GetPvcs(pvc.Name, pvc.Namespace, pvc.Limit, pvc.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取pvc列表成功",
		"data": pvcresp,
	})
}

// 获取pvc详情
func (p *pvc) GetPvcDetail(c *gin.Context) {
	pvc := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
	})
	if err := c.Bind(pvc); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	pvcdetail, err := service.Pvc.GetPvcDetail(pvc.Name, pvc.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取pvc详情成功",
		"data": pvcdetail,
	})
}

// 删除pvc
func (p *pvc) DeletePvc(c *gin.Context) {
	pvc := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
	})
	if err := c.Bind(pvc); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	if err := service.Pvc.DeletePvc(pvc.Name, pvc.Namespace); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除pvc成功",
	})
}

// 更新pvc
func (p *pvc) UpdatePvc(c *gin.Context) {
	pvc := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
	})
	if err := c.ShouldBindJSON(pvc); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	if err := service.Pvc.UpdatePvc(pvc.Namespace, pvc.Content); err != nil {
		c.JSON(400, gin.H{
			"err": "更新pvc失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新pvc成功",
	})
}
