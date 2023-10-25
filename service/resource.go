package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Resources resources

type resources struct {
}
type resourcedata struct {
	Namespace      string `json:"namespace"`
	PodNum         int    `json:"podNum"`
	DeployNum      int    `json:"deployNum"`
	DaemonsetNum   int    `json:"daemonsetNum"`
	StatefulsetNum int    `json:"statefulsetNum"`
	SvcNum         int    `json:"svcNum"`
}

// 获取namespace中的所有资源
func (r *resources) GetAllResource() (resourceNps []*resourcedata, err error) {
	//获取namespace列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace信息失败：" + err.Error())
		return nil, errors.New("获取namespace信息失败：" + err.Error())
	}
	//测试只获取前10个命名空间的数据
	var num = 0
	//获取pod列表
	for _, namespace := range namespaceList.Items {
		if num > 5 {
			break
		}
		//获取pod数量
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取pod信息失败：" + err.Error())
			return nil, errors.New("获取pod信息失败：" + err.Error())
		}
		//获取deployment数量
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取deployment列表失败：" + err.Error())
			return nil, errors.New("获取deployment列表失败：" + err.Error())
		}
		//获取daemonset数量
		daemonsetList, err := K8s.ClientSet.AppsV1().DaemonSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取daemonset列表失败：" + err.Error())
			return nil, errors.New("获取daemonset列表失败：" + err.Error())
		}
		//获取statefulset数量
		statefulsetList, err := K8s.ClientSet.AppsV1().StatefulSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取statefulset列表失败：" + err.Error())
			return nil, errors.New("获取statefulset列表失败：" + err.Error())
		}
		//获取svc数量
		svcList, err := K8s.ClientSet.CoreV1().Services(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取svc列表失败：" + err.Error())
			return nil, errors.New("获取svc列表失败：" + err.Error())
		}
		//组装数据
		resourceNp := &resourcedata{
			Namespace:      namespace.Name,
			PodNum:         len(podList.Items),
			DeployNum:      len(deploymentList.Items),
			DaemonsetNum:   len(daemonsetList.Items),
			StatefulsetNum: len(statefulsetList.Items),
			SvcNum:         len(svcList.Items),
		}
		//添加到podsNps数组中
		resourceNps = append(resourceNps, resourceNp)
		num = num + 1
	}
	return resourceNps, nil
}
