package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct {
}

// 定义获取每个namespace中deployment的数量
type DeploysNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deployNum"`
}

// 定义列表的返回内容，Items是deployment元素列表，Total为deployment元素数量
type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// 定义DeployCreate结构体，用于创建deployment需要的参数属性的定义
type DeployCreate struct {
	Name          string            `json:"name" binding:"required"`
	Namespace     string            `json:"namespace" binding:"required"`
	Img           string            `json:"img" binding:"required"`
	Replicas      int32             `json:"replicas"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Mem           string            `json:"mem"`
	ContainerPort int32             `json:"containerPort"`
	HealthCheck   bool              `json:"healthCheck"`
	HealthPath    string            `json:"healthPath"`
}

// toCells方法用于将deployment类型数组，转换成DataCell类型数组
func (d *deployment) toCell(deployments []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(deployments))
	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成deployment类型数组
func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

// 获取deployment列表，支持过滤、排序、分页
func (d *deployment) GetDeployments(filterName, namespace string, limit, page int) (deploymentResp *DeploymentsResp, err error) {
	//获取deploymentList类型的deployment列表
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取deployemnt列表失败" + err.Error())
		return nil, errors.New("获取deployemnt列表失败" + err.Error())
	}
	//将delpoymentList中的deployment列表（Items），放进dataselector对象中，进行排序
	selectorData := &DataSelector{
		GenericDataList: d.toCell(deploymentList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery:   &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{Limit: limit, Page: page},
		},
	}
	//过滤
	//fmt.Println("过滤之前为：", filterName, selectorData.DataSelectQuery.FilterQuery.Name)
	filtered := selectorData.Filter()
	total := len(filtered.GenericDataList)
	//排序分页
	fmt.Println("传输入的limit: ", filtered.DataSelectQuery.PaginateQuery.Limit)
	data := filtered.Sort().Paginate()
	//将[]DataCell类型的deployment列表转为appsv1.deployment列表
	deployments := d.fromCells(data.GenericDataList)
	return &DeploymentsResp{
		Items: deployments,
		Total: total,
	}, nil
}

// 获取deployment详情
func (d *deployment) GetdeploymentDetail(deploymentName, namespace string) (deploymentDetail *appsv1.Deployment, err error) {
	deploy, err := K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取deployment详情失败" + err.Error())
		return nil, errors.New("获取deployment详情失败" + err.Error())
	}
	return deploy, nil
}

// 删除deployment
func (d *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除deployment失败" + err.Error())
		return errors.New("删除deployment失败" + err.Error())
	}
	return nil
}

// 设置deployment副本数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	//获取autoscalingv1.Scale类型的对象，能点出当前的副本数
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取deployment副本数失败: " + err.Error())
		return 0, errors.New("获取deployment副本数失败: " + err.Error())
	}
	//修改副本数
	scale.Spec.Replicas = int32(scaleNum)
	//设置新的副本数
	newscale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新副本数失败: " + err.Error())
		return 0, errors.New("更新副本数失败: " + err.Error())
	}
	return newscale.Spec.Replicas, nil
}

// 创建deployment,接收DeployCreate对象
func (d *deployment) CreateDeployment(deployData *DeployCreate) (err error) {
	//将deployData中的数据组装成appsv1.Deployment对象
	fmt.Println("************************", *deployData, "资源限制为：", deployData.Cpu, deployData.Mem)
	deployment := &appsv1.Deployment{
		//ObjectMeta中定义资源名、命名空间以及标签
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployData.Name,
			Namespace: deployData.Namespace,
			Labels:    deployData.Label,
		},
		//Spec中定义副本数、选择器、以及pod属性
		Spec: appsv1.DeploymentSpec{
			Replicas: &deployData.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: deployData.Label,
			},
			Template: corev1.PodTemplateSpec{
				//定义pod名和标签
				ObjectMeta: metav1.ObjectMeta{
					Name:   deployData.Name,
					Labels: deployData.Label,
				},
				//定义容器名、镜像和端口
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  deployData.Name,
							Image: deployData.Img,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: deployData.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		//Status定义资源的运行状态，这里由于是新建，传入空的appsv1.DeploymentStatus{}对象即可
		Status: appsv1.DeploymentStatus{},
	}
	//判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
	if deployData.HealthCheck {
		//设置容器的ReadinessProbe
		//若pod中有多个容器，则这里需要使用for循环去定义了
		for _, container := range deployment.Spec.Template.Spec.Containers {
			container.ReadinessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: deployData.HealthPath,
						//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
						//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
						//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
						Port: intstr.IntOrString{
							Type:   0,
							IntVal: deployData.ContainerPort,
						},
					},
				},
				//初始化等待时间
				InitialDelaySeconds: 5,
				//超时时间
				TimeoutSeconds: 5,
				//执行间隔
				PeriodSeconds: 5,
			}
			container.LivenessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: deployData.HealthPath,
						Port: intstr.IntOrString{
							Type:   0,
							IntVal: deployData.ContainerPort,
						},
					},
				},
				InitialDelaySeconds: 15,
				TimeoutSeconds:      5,
				PeriodSeconds:       5,
			}
		}
	}
	//定义容器的limit和request资源
	for _, container := range deployment.Spec.Template.Spec.Containers {
		container.Resources.Limits = map[corev1.ResourceName]resource.Quantity{
			//corev1.ResourceCPU:    resource.MustParse(deployData.Cpu),
			//corev1.ResourceMemory: resource.MustParse(deployData.Mem),
			corev1.ResourceCPU:    resource.MustParse(deployData.Cpu),
			corev1.ResourceMemory: resource.MustParse(deployData.Mem),
		}
		container.Resources.Requests = map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse(deployData.Cpu),
			corev1.ResourceMemory: resource.MustParse(deployData.Mem),
		}
		fmt.Println("赋值完后的资源数据为：", container.Resources.Limits, container.Resources.Requests)
	}
	for _, continerr := range deployment.Spec.Template.Spec.Containers {
		fmt.Println("deployment限制资源：", continerr.Resources.Limits)
		fmt.Println("deployment请求资源：", continerr.Resources.Requests)
	}

	//调用sdk创建deployment
	//_, err = K8s.ClientSet.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	//if err != nil {
	//	logger.Error("创建deployment失败: " + err.Error())
	//	return errors.New("创建deployment失败: " + err.Error())
	//}
	return nil
}

// 重启deployment
func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	//此功能等同于以下kubectl命令
	//kubectl deployment ${service} -p \
	//'{"spec":{"template":{"spec":{"containers":[{"name":"'"${service}"'","env":
	//[{"name":"RESTART_","value":"'$(date +%s)'"}]}]}}}}'

	//使用patchData Map组装数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTART_",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
							}},
						},
					},
				},
			},
		},
	}
	//序列化为字节，因为patch方法只接收字节类型参数
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error("json反序列化失败" + err.Error())
		return errors.New("json反序列化失败" + err.Error())
	}
	//调用patch方法更新deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).
		Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error("更新deployment失败: " + err.Error())
		return errors.New("更新deployment失败: " + err.Error())
	}
	return nil
}

// 更新deployment
func (d *deployment) UpdateDeployment(namespace, content string) (err error) {
	deploy := &appsv1.Deployment{}
	if err := json.Unmarshal([]byte(content), deploy); err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新deployment失败: " + err.Error())
		return errors.New("更新deployment失败: " + err.Error())
	}
	return nil
}

// 获取每个namespace的deployment数量
func (d *deployment) GetDeploymentNumPerNp() (deploysNps []*DeploysNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace列表失败：" + err.Error())
		return nil, errors.New("获取namespace列表失败：" + err.Error())
	}
	var num = 0
	for _, namespace := range namespaceList.Items {
		if num > 10 {
			break
		}
		deployments, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取deployment列表失败：" + err.Error())
			return nil, errors.New("获取deployment列表失败：" + err.Error())
		}
		deploysNp := &DeploysNp{Namespace: namespace.Name, DeployNum: len(deployments.Items)}
		deploysNps = append(deploysNps, deploysNp)
		num = num + 1
	}
	return deploysNps, nil

}
