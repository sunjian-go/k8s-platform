package service

import (
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

func (k *k8s) Init() {
	//1.将kubeconfig格式化为rest.config类型
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		panic("创建k8s配置失败，" + err.Error())
	}

	//2.通过config创建clientset
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建k8s clientSet失败, " + err.Error())
	} else {
		logger.Info("创建k8s clinetSet成功")
	}
	k.ClientSet = clientset
}
