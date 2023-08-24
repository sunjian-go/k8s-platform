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
	//workflowResp, err := dao.Workflow.GetWorkflow("nginx", "default", 1, 1)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for _, data := range workflowResp.Items {
	//	fmt.Println("查询出来的workflow数据为：", data)
	//}
	//data, _ := dao.Workflow.GetById(1)
	//fmt.Println("按条件查询出来的数据为：", data)
	//初始化路由
	r := gin.Default()
	r.Use(middle.Cors()) //放在初始化路由之前
	controller.Router.InitApiRouter(r)

	r.Run(config.ListenAddr)
}
