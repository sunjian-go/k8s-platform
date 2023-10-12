package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	Networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct {
}
type IngResp struct {
	Items []Networkingv1.Ingress `json:"items"`
	Total int                    `json:"total"`
}
type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
}
type HttpPath struct {
	Path        string                `json:"path"`
	PathType    Networkingv1.PathType `json:"path_type"`
	ServiceName string                `json:"service_name"`
	ServicePort int32                 `json:"service_port"`
}

func (i *ingress) toCell(ingresses []Networkingv1.Ingress) []DataCell {
	cells := make([]DataCell, len(ingresses))
	for i := range ingresses {
		cells[i] = ingressCell(ingresses[i])
	}
	return cells
}

func (i *ingress) fromCells(cells []DataCell) []Networkingv1.Ingress {
	ingresses := make([]Networkingv1.Ingress, len(cells))
	for i := range cells {
		ingresses[i] = Networkingv1.Ingress(cells[i].(ingressCell))
	}
	return ingresses
}

// 获取ing列表
func (i *ingress) GetIngresses(ingName, namespace string, limit, page int) (ingResp *IngResp, err error) {
	ingList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取ingress列表失败：" + err.Error())
		return nil, errors.New("获取ingress列表失败：" + err.Error())
	}

	data := &DataSelector{
		GenericDataList: i.toCell(ingList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: ingName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	newdata := data.Filter()
	total := len(newdata.GenericDataList)
	ings := i.fromCells(newdata.Sort().Paginate().GenericDataList)
	fmt.Println("获取到ing列表为：", ings)
	return &IngResp{
		Items: ings,
		Total: total,
	}, nil
}

// 获取Ingress详情
func (i *ingress) GetIngDetail(ingName, namespace string) (ing *Networkingv1.Ingress, err error) {
	ingress, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取ingress详情失败: " + err.Error())
		return nil, errors.New("获取ingress详情失败: " + err.Error())
	}
	return ingress, nil
}

// 创建Ingress
func (i *ingress) CreateIngress(data *IngressCreate) (err error) {
	//声明Networkingv1.IngressRule和Networkingv1.HTTPIngressPath变量，后面组装数据用到
	var ingressRules []Networkingv1.IngressRule
	var httpIngressPATHs []Networkingv1.HTTPIngressPath
	//将data中的数据组装成Networkingv1.Ingress对象
	ingress := &Networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Status: Networkingv1.IngressStatus{},
	}
	//第一层for循环是将host组装成Networkingv1.IngressRule类型的对象
	// 一个host对应一个ingressrule，每个ingressrule中包含一个host和多个path
	for key, value := range data.Hosts {
		ir := Networkingv1.IngressRule{
			Host: key,
			//这里先将Networkingv1.HTTPIngressRuleValue类型中的Paths置为空，后面组装好数据再赋值
			IngressRuleValue: Networkingv1.IngressRuleValue{
				HTTP: &Networkingv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _, httppath := range value {
			hip := Networkingv1.HTTPIngressPath{
				Path:     httppath.Path,
				PathType: &httppath.PathType,
				Backend: Networkingv1.IngressBackend{
					Service: &Networkingv1.IngressServiceBackend{
						Name: httppath.ServiceName,
						Port: Networkingv1.ServiceBackendPort{
							Number: httppath.ServicePort,
						},
					},
				},
			}
			//将每个hip对象组装成数组
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}
		//给Paths赋值，前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		//将每个ir对象组装成数组，这个ir对象就是IngressRule，每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}
	//将ingressRules对象加入到ingress的规则中
	ingress.Spec.Rules = ingressRules
	//创建ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建ingress失败: " + err.Error())
		return errors.New("创建ingress失败: " + err.Error())
	}
	return nil
}

// 删除Ingress
func (i *ingress) DeleteIngress(ingName, namespace string) (err error) {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除ingress失败：" + err.Error())
		return errors.New("删除ingress失败：" + err.Error())
	}
	return nil
}

// 更新Ingress
func (i *ingress) UpdateIngress(namespace, content string) (err error) {
	ing := new(Networkingv1.Ingress)
	if err := json.Unmarshal([]byte(content), ing); err != nil {
		logger.Error("解码失败: " + err.Error())
		return errors.New("解码失败: " + err.Error())
	}
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ing, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新ingress失败：" + err.Error())
		return errors.New("更新ingress失败：" + err.Error())
	}
	return nil
}
