package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
)

var Pv pv

type pv struct {
}

// 获取pv列表
func (p *pv) GetPvs(c *gin.Context) {
	pvResp := new(struct {
		Name  string `form:"name"`
		Limit int    `form:"limit"`
		Page  int    `form:"page"`
	})
	if err := c.Bind(pvResp); err != nil {
		logger.Info("绑定数据失败: ", err.Error())
		c.JSON(400, gin.H{
			"err": "绑定数据失败: " + err.Error(),
		})
		return
	}
	pvs, err := service.Pv.GetPvs(pvResp.Name, pvResp.Limit, pvResp.Page)
	if err != nil {
		logger.Info(err.Error())
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
