package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Secret secret

type secret struct {
}

type secretResp struct {
	Items []corev1.Secret
	Total int
}

func (s *secret) toCell(secrets []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(secrets))
	for i := range secrets {
		cells[i] = secretCell(secrets[i])
	}
	return cells
}

func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return secrets
}

// 获取secret列表
func (s *secret) GetSecrets(secretName, namespace string, limit, page int) (secretresp *secretResp, err error) {
	secrets, err := K8s.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取secret列表失败: " + err.Error())
		return nil, errors.New("获取secret列表失败: " + err.Error())
	}
	data := &DataSelector{
		GenericDataList: s.toCell(secrets.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: secretName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	newdata := data.Filter()
	total := len(newdata.GenericDataList)
	secretList := s.fromCells(newdata.Sort().Paginate().GenericDataList)
	return &secretResp{
		Items: secretList,
		Total: total,
	}, nil
}

// 获取secret详情
func (s *secret) GetSecretDetail(secretName, namespace string) (secret *corev1.Secret, err error) {
	secret, err = K8s.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取secret详情失败: " + err.Error())
		return nil, errors.New("获取secret详情失败: " + err.Error())
	}
	return secret, nil
}

// 删除secret
func (s *secret) DeleteSecret(secretName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除secret失败: " + err.Error())
		return errors.New("删除secret失败: " + err.Error())
	}
	return nil
}

// 更新secret
func (s *secret) UpdateSecret(namespace, content string) (err error) {
	secret := new(corev1.Secret)
	if err := json.Unmarshal([]byte(content), secret); err != nil {
		logger.Error("解码失败: " + err.Error())
		return errors.New("解码失败: " + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新secret失败: " + err.Error())
		return errors.New("更新secret失败: " + err.Error())
	}
	return nil
}
