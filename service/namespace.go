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

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (n *namespace) toCell(namespace []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(namespace))
	for i := range namespace {
		cells[i] = nsCell(namespace[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	ns := make([]corev1.Namespace, len(cells))
	for i := range cells {
		ns[i] = corev1.Namespace(cells[i].(nsCell))
	}
	return ns
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

// /按需获取namespace列表，支持过滤、排序、分页，专攻前端namespace页面使用
func (n *namespace) GetNamespaceList(filterNamespace string, limit, page int) (namespaceresp *namespaceResp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败：" + err.Error())
		return nil, errors.New("获取namespace列表失败：" + err.Error())
	}
	//实例化dataSelector对象
	selecttableData := &DataSelector{
		GenericDataList: n.toCell(namespaceList.Items), //将pods列表转换为DataCell类型赋值
		DataSelectQuery: &DataSelectQuery{ //
			FilterQuery: &FilterQuery{
				Name: filterNamespace, //将传进来的需要查找的Name赋值给该结构体
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit, //将传进来的页数和每页的数量赋值
				Page:  page,
			},
		},
	}
	//先过滤
	filterd := selecttableData.Filter()
	total := len(filterd.GenericDataList)   //计算过滤好的目标pod列表的长度
	namespaces := filterd.Sort().Paginate() //连续调用排序和分页方法
	//将[]DataCell类型的pod列表转为v1.pod列表
	namespacev1 := n.fromCells(namespaces.GenericDataList)
	return &namespaceResp{
		Namespaces: namespacev1,
		Total:      total,
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
