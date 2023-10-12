package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/dao"
	"k8s-platform/service"
)

var Styles styles

type styles struct {
}

// 获取颜色
func (s *styles) GetColor(c *gin.Context) {
	colors, err := service.Styles.GetColor()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取颜色信息成功",
		"data": colors,
	})
}

// 更新颜色
func (s *styles) UpdateColor(c *gin.Context) {
	colors := new(dao.ColorResp)
	if err := c.ShouldBind(colors); err != nil {
		fmt.Println("绑定数据失败：", err.Error())
		c.JSON(400, gin.H{
			"err": "绑定数据失败：" + err.Error(),
		})
		return
	}
	fmt.Println("要更新的颜色为：", colors)
	err := service.Styles.UpdateColor(colors.Background, colors.Color)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"msg": "更新颜色成功",
	})
}
