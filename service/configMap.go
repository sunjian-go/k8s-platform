package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ConfigMap configMap

type configMap struct {
}
type configMapResp struct {
	Items []corev1.ConfigMap `json:"items"`
	Total int                `json:"total"`
}

func (c *configMap) toCell(configMaps []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(configMaps))
	for i := range configMaps {
		cells[i] = cmCell(configMaps[i])
	}
	return cells
}

func (c *configMap) fromCells(cells []DataCell) []corev1.ConfigMap {
	configMaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		configMaps[i] = corev1.ConfigMap(cells[i].(cmCell))
	}
	return configMaps
}

// 获取cm列表
func (c *configMap) GetConfigMaps(cmName, namespace string, limit, page int) (configmapResp *configMapResp, err error) {
	cms, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取cm列表失败: " + err.Error())
		return nil, errors.New("获取cm列表失败: " + err.Error())
	}
	data := &DataSelector{
		GenericDataList: c.toCell(cms.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: cmName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	newdata := data.Filter()
	total := len(newdata.GenericDataList)
	cmList := c.fromCells(newdata.Sort().Paginate().GenericDataList)
	return &configMapResp{
		Items: cmList,
		Total: total,
	}, nil
}

// 获取cm详情
func (c *configMap) GetConfigMapDetail(cmName, namespace string) (configmap *corev1.ConfigMap, err error) {
	configmap, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cmName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取cm详情失败: " + err.Error())
		return nil, errors.New("获取cm详情失败: " + err.Error())
	}
	return configmap, nil
}

// 删除cm
func (c *configMap) DeleteConfigMap(cmName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), cmName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除cm失败: " + err.Error())
		return errors.New("删除cm失败: " + err.Error())
	}
	return nil
}

// 更新cm
func (c *configMap) UpdateConfigMap(namespace, content string) (err error) {
	configmap := new(corev1.ConfigMap)
	if err := json.Unmarshal([]byte(content), configmap); err != nil {
		logger.Error("解码失败: " + err.Error())
		return errors.New("解码失败: " + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configmap, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新cm失败: " + err.Error())
		return errors.New("更新cm失败: " + err.Error())
	}
	return nil
}
