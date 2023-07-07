package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DaemonSet daemonSet

type daemonSet struct {
}
type DaemonSetResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

// toCells方法用于将daemonSet类型数组，转换成DataCell类型数组
func (d *daemonSet) toCell(daemonSets []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(daemonSets))
	for i := range daemonSets {
		cells[i] = daemonSetCell(daemonSets[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成daemonSetCell类型数组
func (d *daemonSet) fromCells(cells []DataCell) []appsv1.DaemonSet {
	daemonSets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		//cells[i].(daemonSetCell)就使用到了断言,断言后转换成了daemonSetCell类型，然后又转换成了daemonSet类型
		daemonSets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return daemonSets
}

// 获取daemonset列表,支持过滤、分页、排序
func (d *daemonSet) GetDaemonSets(daemonSetName, namespace string, limit, page int) (daemonSets *DaemonSetResp, err error) {
	daemonsets, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取daemonSet列表失败： " + err.Error())
		return nil, errors.New("获取daemonSet列表失败： " + err.Error())
	}
	//组装好数据准备下一步
	daemonset := &DataSelector{
		GenericDataList: d.toCell(daemonsets.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: daemonSetName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//过滤
	filter := daemonset.Filter()
	total := len(filter.GenericDataList)
	//排序、分页
	data := filter.Sort().Paginate()
	daemonsetss := d.fromCells(data.GenericDataList)
	return &DaemonSetResp{
		Items: daemonsetss,
		Total: total,
	}, nil
}

// 获取daemonset详情
func (d *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (daemonsetDetail *appsv1.DaemonSet, err error) {
	daemonsetDetail, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取daemonset详情失败：" + err.Error())
		return nil, errors.New("获取daemonset详情失败：" + err.Error())
	}
	return daemonsetDetail, nil
}

// 删除daemonset
func (d *daemonSet) DeleteDaemonSet(daemonsetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonsetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除daemonset失败：" + err.Error())
		return errors.New("删除daemonset失败：" + err.Error())
	}
	return nil
}

// 更新daemonset
func (d *daemonSet) UpdateDaemonSet(namespace, content string) (err error) {
	daemonset := new(appsv1.DaemonSet)
	//将传进来的字符串，转为字符切片类型反序列化到定义好的appsv1.DaemonSet结构体中
	if err := json.Unmarshal([]byte(content), daemonset); err != nil {
		logger.Error("反序列化失败：" + err.Error())
		return errors.New("反序列化失败：" + err.Error())
	}
	_, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonset, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新daemonset失败: " + err.Error())
		return errors.New("更新daemonset失败: " + err.Error())
	}
	return nil
}
