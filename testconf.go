package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

func main() {
	currentPath, _ := os.Getwd()
	confPath := currentPath + "/conf/kube_conf.ini"
	_, err := os.Stat(confPath)
	if err != nil {
		//panic(errors.New(fmt.Sprintf("file is not found %s", confPath)))
		panic("file is not found " + confPath)
	}
	// 加载配置
	config, err := goconfig.LoadConfigFile(confPath)
	if err != nil {
		log.Fatal("读取配置文件出错:", err)
	}
	// 获取 section
	kubeconf, _ := config.GetSection("kube")
	fmt.Println("配置文件内容：", kubeconf)
	if kubeconf["LogMode"] == "true" {
		fmt.Println("LogMode：", kubeconf["LogMode"])
	} else {
		fmt.Println("LogMode：", kubeconf["LogMode"])
	}
	fmt.Printf("类型为：%T\n", kubeconf["MaxIdleConns"])

}
