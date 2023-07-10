package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namespace namespace

type namespace struct {
}
type namespaceResp struct {
	Namespaces []corev1.Namespace `json:"namespaces"`
	Total      int                `json:"total"`
}

// 获取namespace列表
func (n *namespace) GetNamespaces() (namespaces *namespaceResp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败：" + err.Error())
		return nil, errors.New("获取namespace列表失败：" + err.Error())
	}
	return &namespaceResp{
		Namespaces: namespaceList.Items,
		Total:      len(namespaceList.Items),
	}, nil
}

// 获取Namespace详情
func (n *namespace) GetNamespaceDetail(namespaceName string) (namespace *corev1.Namespace, err error) {
	namespace, err = K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取namespace详情失败: " + err.Error())
		return nil, errors.New("获取namespace详情失败: " + err.Error())
	}
	return namespace, nil
}

// 删除Namespace
func (n *namespace) DeleteNamespace(namespaceName string) (err error) {
	err = K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除namespace失败: " + err.Error())
		return errors.New("删除namespace失败: " + err.Error())
	}
	return nil
}
