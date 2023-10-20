package controller

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/service"
)

var Secret secret

type secret struct {
}

// 获取secret列表
func (s *secret) GetSecrets(c *gin.Context) {
	secret := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Limit     int    `form:"limit"`
		Page      int    `form:"page"`
	})
	if err := c.Bind(secret); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	secrets, err := service.Secret.GetSecrets(secret.Name, secret.Namespace, secret.Limit, secret.Page)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取secret列表成功",
		"data": secrets,
	})
}

// 获取secret详情
func (s *secret) GetSecretDetail(c *gin.Context) {
	secret := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(secret); err != nil {
		c.JSON(400, gin.H{
			"err":  "绑定数据失败",
			"data": nil,
		})
		return
	}
	secretdetail, err := service.Secret.GetSecretDetail(secret.Name, secret.Namespace)
	if err != nil {
		c.JSON(400, gin.H{
			"err":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取secret详情成功",
		"data": secretdetail,
	})
}

// 删除secret
func (s *secret) DeleteSecret(c *gin.Context) {
	secret := new(struct {
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
	})
	if err := c.Bind(secret); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
		return
	}
	if err := service.Secret.DeleteSecret(secret.Name, secret.Namespace); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "删除secret成功",
	})
}

// 更新secret
func (s *secret) UpdateSecret(c *gin.Context) {
	secret := new(struct {
		Namespace string `json:"namespace" binding:"required"`
		Content   string `json:"content" binding:"required"`
	})
	if err := c.ShouldBindJSON(secret); err != nil {
		c.JSON(400, gin.H{
			"err": "绑定数据失败",
		})
	}
	if err := service.Secret.UpdateSecret(secret.Namespace, secret.Content); err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "更新secret成功",
	})
}
