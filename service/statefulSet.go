package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var StatefulSet statefulSet

type statefulSet struct {
}
type statefulSetResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

// toCells方法用于将statefulSetCell类型数组，转换成DataCell类型数组
func (s *statefulSet) toCell(statefulSets []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(statefulSets))
	for i := range statefulSets {
		cells[i] = statefulSetCell(statefulSets[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成statefulSetCell类型数组
func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulSets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		//cells[i].(statefulSetCell)就使用到了断言,断言后转换成了statefulSetCell类型，然后又转换成了statefulSet类型
		statefulSets[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}
	return statefulSets
}

// 获取StatefulSet列表,支持过滤、分页、排序
func (s *statefulSet) GetStatefulSets(statefulSetName, namespace string, limit, page int) (statefulSetresp *statefulSetResp, err error) {
	statefulsets, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取statefulset列表失败: " + err.Error())
		return nil, errors.New("获取statefulset列表失败: " + err.Error())
	}
	dataselector := &DataSelector{
		GenericDataList: s.toCell(statefulsets.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: statefulSetName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//过滤
	filter := dataselector.Filter()
	total := len(filter.GenericDataList)
	//排序、分页
	data := filter.Sort().Paginate()
	statefulsetData := s.fromCells(data.GenericDataList)
	return &statefulSetResp{
		Items: statefulsetData,
		Total: total,
	}, nil
}

// 获取StatefulSet详情
func (s *statefulSet) GetStatefulSetDetail(statefuleSetName, namespace string) (statefulset *appsv1.StatefulSet, err error) {
	statefulset, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefuleSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取statefulSet详情失败: " + err.Error())
		return nil, errors.New("获取statefulSet详情失败: " + err.Error())
	}
	return statefulset, nil
}

// 删除StatefulSet
func (s *statefulSet) DeleteStatefulSet(statefulSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除statefulSet失败: " + err.Error())
		return errors.New("删除statefulSet失败: " + err.Error())
	}
	return nil
}

// 更新StatefulSet
func (s *statefulSet) UpdateStatefulSet(namespace, content string) (err error) {
	//先解码
	statefulset := new(appsv1.StatefulSet)
	if err := json.Unmarshal([]byte(content), statefulset); err != nil {
		logger.Error("解码失败: " + err.Error())
		return errors.New("解码失败: " + err.Error())
	}
	//再更新
	_, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulset, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新statefulSet失败: " + err.Error())
		return errors.New("更新statefulSet失败: " + err.Error())
	}
	return nil
}
