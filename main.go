package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/service"
)

func main() {
	//初始化k8s client
	service.K8s.Init()
	//初始化路由
	r := gin.Default()
	controller.Router.InitApiRouter(r)
	r.Run(config.ListenAddr)
}
