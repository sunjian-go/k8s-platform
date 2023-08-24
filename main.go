package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/middle"
	"k8s-platform/service"
)

func main() {
	//初始化k8s client
	service.K8s.Init()
	//初始化mysql
	db.Init()
	//初始化路由
	r := gin.Default()
	r.Use(middle.JWTAuth()) //加载jwt中间件，用于token验证
	r.Use(middle.Cors())    //加载跨域中间件（放在初始化路由之前）
	controller.Router.InitApiRouter(r)
	r.Run(config.ListenAddr)
}
