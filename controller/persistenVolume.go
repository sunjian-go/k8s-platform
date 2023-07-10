package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var Pv pv

type pv struct {
}

// 获取pv列表
func (p *pv) GetPvs(c *gin.Context) {
	pvs, err := service.Pv.GetPvs()
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取pv列表成功",
		"data": pvs,
	})
}

// 获取Pv详情
func (p *pv) GetPvDetail(c *gin.Context) {
	//GET请求
	pvname := c.Query("name")
	pv, err := service.Pv.GetPvDetail(pvname)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取pv详情成功",
		"data": pv,
	})
}

// 删除Pv
func (p *pv) DeletePv(c *gin.Context) {
	pvname := c.Query("name")
	if err := service.Pv.DeletePv(pvname); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除pv成功",
	})
}
