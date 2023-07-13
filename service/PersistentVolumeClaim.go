package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pvc pvc

type pvc struct {
}

type pvcResp struct {
	Items []corev1.PersistentVolumeClaim `json:"items"`
	Total int                            `json:"total"`
}

func (p *pvc) toCell(pvcs []corev1.PersistentVolumeClaim) []DataCell {
	cells := make([]DataCell, len(pvcs))
	for i := range pvcs {
		cells[i] = pvcCell(pvcs[i])
	}
	return cells
}

func (p *pvc) fromCells(cells []DataCell) []corev1.PersistentVolumeClaim {
	pvcs := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		pvcs[i] = corev1.PersistentVolumeClaim(cells[i].(pvcCell))
	}
	return pvcs
}

// 获取pvc列表
func (p *pvc) GetPvcs(pvcName, namespace string, limit, page int) (pvcresp *pvcResp, err error) {
	pvcs, err := K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取pvc列表失败: " + err.Error())
		return nil, errors.New("获取pvc列表失败: " + err.Error())
	}
	data := &DataSelector{
		GenericDataList: p.toCell(pvcs.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: pvcName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	newdata := data.Filter()
	total := len(newdata.GenericDataList)
	pvcList := p.fromCells(newdata.Sort().Paginate().GenericDataList)
	return &pvcResp{
		Items: pvcList,
		Total: total,
	}, nil
}

// 获取pvc详情
func (p *pvc) GetPvcDetail(pvcName, namespace string) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvc, err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取pvc详情失败: " + err.Error())
		return nil, errors.New("获取pvc详情失败: " + err.Error())
	}
	return pvc, nil
}

// 删除pvc
func (p *pvc) DeletePvc(pvcName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除pvc失败：" + err.Error())
		return errors.New("删除pvc失败：" + err.Error())
	}
	return nil
}

// 更新pvc
func (p *pvc) UpdatePvc(namespace, content string) (err error) {
	pvc := new(corev1.PersistentVolumeClaim)
	if err := json.Unmarshal([]byte(content), pvc); err != nil {
		logger.Error("解码失败：" + err.Error())
		return errors.New("解码失败：" + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新pvc失败：" + err.Error())
		return errors.New("更新pvc失败：" + err.Error())
	}
	return nil
}
