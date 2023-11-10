package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/dao"
	"k8s-platform/db"
	"k8s-platform/middle"
	"k8s-platform/service"
	"net/http"
)

func main() {
	//初始化k8s client
	service.K8s.Init()
	//初始化mysql
	db.Init()
	colorsResp, _ := dao.Styles.GetColor()
	if colorsResp.Background == "" || colorsResp.Color == "" {
		//如果styles表里没数据就进行初始化
		_ = dao.Styles.InitStyles()
		fmt.Println("初始化styles表成功")
	}
	//_ = dao.Styles.UpdateColor("red", "green")
	//初始化路由
	r := gin.Default()
	r.Use(middle.Cors())    //加载跨域中间件(一定要先跨域，再加载jwt)
	r.Use(middle.JWTAuth()) //加载jwt中间件，用于token验证
	controller.Router.InitApiRouter(r)
	//启动websocket
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":8081", nil)
		fmt.Println("ws服务已启动。。。")
	}()
	r.Run(config.ListenAddr)
	//关闭数据库连接
	db.Close()
}
