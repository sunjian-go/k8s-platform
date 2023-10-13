package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var SVC service

type service struct {
}
type serviceResp struct {
	Services []corev1.Service `json:"services"`
	Total    int              `json:"total"`
}

//定义ServiceCreate结构体，用于创建service需要的参数属性的定义

type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	Label         map[string]string `json:"label"`
}

// toCells方法用于将service类型数组，转换成DataCell类型数组
func (s *service) toCell(services []corev1.Service) []DataCell {
	cells := make([]DataCell, len(services))
	for i := range services {
		cells[i] = serviceCell(services[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成service类型数组
func (s *service) fromCells(cells []DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i := range cells {
		//cells[i].(service)就使用到了断言,断言后转换成了serviceCell类型，然后又转换成了service类型
		services[i] = corev1.Service(cells[i].(serviceCell))
	}
	return services
}

// 获取svc列表,支持过滤、排序、分页
func (s *service) GetSvcs(svcname, namespace string, limit, page int) (serviceresp *serviceResp, err error) {
	//先获取svc列表
	svcList, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取svc列表失败: " + err.Error())
		return nil, errors.New("获取svc列表失败: " + err.Error())
	}
	//组装dataselector进行过滤、排序、分页
	data := &DataSelector{
		GenericDataList: s.toCell(svcList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: svcname,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	//先过滤
	newdata := data.Filter()
	//计算过滤出来的元素切片的长度
	total := len(newdata.GenericDataList)
	//排序、分页
	svcs := s.fromCells(newdata.Sort().Paginate().GenericDataList)
	//组装数据返回
	return &serviceResp{
		Services: svcs,
		Total:    total,
	}, nil
}

// 获取Service详情
func (s *service) GetSvcDetail(svcname, namespace string) (svc *corev1.Service, err error) {
	svc, err = K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), svcname, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取svc详情失败: " + err.Error())
		return nil, errors.New("获取svc详情失败: " + err.Error())
	}
	return svc, nil
}

// 删除Service
func (s *service) DeleteSvc(svcname, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), svcname, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除svc失败: " + err.Error())
		return errors.New("删除svc失败: " + err.Error())
	}
	return nil
}

// 创建Service
func (s *service) CreateSvc(data *ServiceCreate) (err error) {
	svc := &corev1.Service{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		//Spec中定义类型，端口，选择器
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
	}
	//默认ClusterIP,这里是判断NodePort,添加配置
	if data.NodePort != 0 && data.Type == "NodePort" {
		svc.Spec.Ports[0].NodePort = data.NodePort
	}
	//创建svc
	_, err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建svc失败: " + err.Error())
		return errors.New("创建svc失败: " + err.Error())
	}
	return nil
}

// 更新Service
func (s *service) UpdateSvc(namespace, content string) (err error) {
	//创建一个svc类型的原生结构体
	svcontent := new(corev1.Service)
	//解码到svc结构体
	if err := json.Unmarshal([]byte(content), svcontent); err != nil {
		logger.Error("解码失败: " + err.Error())
		return errors.New("解码失败: " + err.Error())
	}
	//更新
	_, err = K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), svcontent, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新svc失败: " + err.Error())
		return errors.New("更新svc失败: " + err.Error())
	}
	return nil
}
