package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Node node

type node struct {
}
type NodeResp struct {
	Items []corev1.Node `json:"items"`
	Total int           `json:"total"`
}

// 获取node列表
func (n *node) GetNodes() (noderesp *NodeResp, err error) {
	nodeList, err := K8s.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取node列表失败: " + err.Error())
		return nil, errors.New("获取node列表失败: " + err.Error())
	}

	return &NodeResp{
		Items: nodeList.Items,
		Total: len(nodeList.Items),
	}, nil
}

// 获取node详情
func (n *node) GetNodeDetail(nodeName string) (node *corev1.Node, err error) {
	node, err = K8s.ClientSet.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取node详情失败: " + err.Error())
		return nil, errors.New("获取node详情失败: " + err.Error())
	}
	return node, nil
}
