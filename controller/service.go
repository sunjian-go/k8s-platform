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
		Name      string `form:"name" binding:"required"`
		Namespace string `form:"namespace" binding:"required"`
		Limit     int    `form:"limit" binding:"required"`
		Page      int    `form:"page" binding:"required"`
	})
	service.SVC.GetSvcs()
}

// 获取Service详情
func (s *svc) GetSvcDetail(c *gin.Context) {
	//GET请求
}

// 创建Service
func (s *svc) CreateSvc(c *gin.Context) {
	//POST请求
}

// 删除Service
func (s *svc) DeleteSvc(c *gin.Context) {
	//DELETE请求
}

// 更新Service
func (s *svc) UpdateSvc(c *gin.Context) {
	//PUT请求
}
