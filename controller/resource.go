package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
)

var Resources resources

type resources struct {
}

// 获取namespace中所有资源
func (r *resources) GetAllResources(c *gin.Context) {
	resourceNps, err := service.Resources.GetAllResource()
	if err != nil {
		logger.Error("获取资源失败：" + err.Error())
		c.JSON(400, gin.H{
			"err":  "获取资源失败：" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取资源成功",
		"data": resourceNps,
	})
}
