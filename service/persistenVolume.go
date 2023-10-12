package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pv pv

type pv struct {
}
type pvResp struct {
	Pvs   []corev1.PersistentVolume `json:"pvs"`
	Total int                       `json:"total"`
}

// 获取pv列表
//
//	func (p *pv) GetPvs() (pvresp *pvResp, err error) {
//		pvList, err := K8s.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
//		if err != nil {
//			logger.Error("获取pv列表失败: " + err.Error())
//			return nil, errors.New("获取pv列表失败: " + err.Error())
//		}
//		return &pvResp{
//			Pvs:   pvList.Items,
//			Total: len(pvList.Items),
//		}, nil
//	}
func (p *pv) toCell(pvs []corev1.PersistentVolume) []DataCell {
	cells := make([]DataCell, len(pvs))
	for i := range pvs {
		cells[i] = pvCell(pvs[i])
	}
	return cells
}

func (p *pv) fromCells(cells []DataCell) []corev1.PersistentVolume {
	pvs := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		pvs[i] = corev1.PersistentVolume(cells[i].(pvCell))
	}
	return pvs
}
func (p *pv) GetPvs(pvName string, limit, page int) (pvresp *pvResp, err error) {
	pvList, err := K8s.ClientSet.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取pv列表失败: " + err.Error())
		return nil, errors.New("获取pv列表失败: " + err.Error())
	}
	data := &DataSelector{
		GenericDataList: p.toCell(pvList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: pvName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	newdata := data.Filter()
	total := len(newdata.GenericDataList)
	pvs := p.fromCells(newdata.Sort().Paginate().GenericDataList)

	return &pvResp{
		Pvs:   pvs,
		Total: total,
	}, nil
}

// 获取Pv详情
func (p *pv) GetPvDetail(pvname string) (pv *corev1.PersistentVolume, err error) {
	pv, err = K8s.ClientSet.CoreV1().PersistentVolumes().Get(context.TODO(), pvname, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取pv详情失败: " + err.Error())
		return nil, errors.New("获取pv详情失败: " + err.Error())
	}
	return pv, nil
}

// 删除Pv
func (p *pv) DeletePv(pvname string) (err error) {
	err = K8s.ClientSet.CoreV1().PersistentVolumes().Delete(context.TODO(), pvname, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除pv失败: " + err.Error())
		return errors.New("删除pv失败: " + err.Error())
	}
	return nil
}
