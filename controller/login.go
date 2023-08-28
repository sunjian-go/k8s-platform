package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var Login login

type login struct {
}

func (l *login) Login(c *gin.Context) {
	user := new(service.User)
	if err := c.ShouldBind(user); err != nil {
		c.JSON(400, gin.H{
			"msg": "数据绑定失败：" + err.Error(),
		})
		return
	}
	if err := service.Login.Login(user); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "登录成功",
	})
}
