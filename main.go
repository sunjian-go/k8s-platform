package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/dao"
	"k8s-platform/db"
	"k8s-platform/service"
)

func main() {
	//初始化k8s client
	service.K8s.Init()
	//初始化mysql
	db.Init()
	workflowResp, _ := dao.Workflow.GetWorkflow("nginx", 1, 1)
	for _, data := range workflowResp.Items {
		fmt.Println("查询出来的数据为：", data)
	}

	//初始化路由
	r := gin.Default()
	controller.Router.InitApiRouter(r)
	r.Run(config.ListenAddr)
}
