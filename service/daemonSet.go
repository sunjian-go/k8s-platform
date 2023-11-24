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
)

var DaemonSet daemonSet

type daemonSet struct {
}
type DaemonSetResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

// 定义DaemonSetCreate结构体，用于创建DaemonSet需要的参数属性的定义
type DaemonSetCreate struct {
	Name              string            `json:"name" binding:"required"`
	Namespace         string            `json:"namespace" binding:"required"`
	Label             map[string]string `json:"label"`
	Cpu               string            `json:"cpu"`
	Mem               string            `json:"mem"`
	HealthCheck       bool              `json:"healthCheck"`
	HealthPath        string            `json:"healthPath"`
	Volume            []Volumes         `json:"volume"`
	NodeSelectorLabel map[string]string `json:"nodeSelectorLabel"`
	Containers        []Container       `json:"containers"`
}

// 定义卷结构体
type Volumes struct {
	VolumeName string `json:"volumeName"`
	Type       string `json:"type"`
	Context    string `json:"context"`
}

// 定义卷挂载结构体
type MontVolumes struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly"`
	SubPath   string `json:"subPath"`
}

// 环境变量
type Env struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
}

// 定义容器组结构体
type Container struct {
	Name            string            `json:"name"`
	Image           string            `json:"image"`
	Ports           []*ContainerPorts `json:"ports"`
	MontVolume      []*MontVolumes    `json:"montVolume"`
	Envs            []*Env            `json:"envs"`
	ImagePullpolicy string            `json:"imagePullpolicy"`
}

// 定义容器端口组结构体
type ContainerPorts struct {
	PortName      string `json:"portName"`
	ContainerPort int32  `json:"containerPort"`
	HostPort      int32  `json:"hostPort"`
	HostIP        string `json:"hostIP"`
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

// 创建daemonset
func (d *daemonSet) CreateDaemonSet(daemonsetData *DaemonSetCreate) (err error) {
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonsetData.Name,
			Namespace: daemonsetData.Namespace,
			Labels:    daemonsetData.Label,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: daemonsetData.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   daemonsetData.Name,
					Labels: daemonsetData.Label,
				},
			},
		},
		Status: appsv1.DaemonSetStatus{},
	}
	//判断是否有卷需要挂载
	if daemonsetData.Volume != nil {
		var volumeSource corev1.VolumeSource
		volumes := make([]corev1.Volume, len(daemonsetData.Volume))
		for i, _ := range daemonsetData.Volume {
			switch daemonsetData.Volume[i].Type {
			case "configMap":
				volumeSource = corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: daemonsetData.Volume[i].Context,
						},
					},
				}
			case "HostPath":
				volumeSource = corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: daemonsetData.Volume[i].Context,
					},
				}
			case "EmptyDir":
				volumeSource = corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				}
			case "PersistentVolumeClaim":
				volumeSource = corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: daemonsetData.Volume[i].Context,
					},
				}
			case "Secret":
				volumeSource = corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: daemonsetData.Volume[i].Context,
					},
				}
			}
			//给volume数组赋值
			volumes[i] = corev1.Volume{
				Name:         daemonsetData.Volume[i].VolumeName,
				VolumeSource: volumeSource,
			}
			fmt.Println("赋值前：", volumes)
		}
		daemonset.Spec.Template.Spec.Volumes = volumes
		fmt.Println("卷数据为：", daemonset.Spec.Template.Spec.Volumes)
	}
	//判断是否使用节点亲和性
	if daemonsetData.NodeSelectorLabel != nil {
		daemonset.Spec.Template.Spec.NodeSelector = daemonsetData.NodeSelectorLabel
	}
	//组装每个容器需要的端口配置
	containers := make([]corev1.Container, len(daemonsetData.Containers))
	for i, _ := range daemonsetData.Containers {
		containers[i] = corev1.Container{
			Name:  daemonsetData.Containers[i].Name,
			Image: daemonsetData.Containers[i].Image,
		}
		//组装每个容器需要的端口组
		ports := make([]corev1.ContainerPort, len(daemonsetData.Containers[i].Ports))
		for j, _ := range daemonsetData.Containers[i].Ports {
			ports[j] = corev1.ContainerPort{
				Name:          daemonsetData.Containers[i].Ports[j].PortName,
				ContainerPort: daemonsetData.Containers[i].Ports[j].ContainerPort,
				Protocol:      corev1.ProtocolTCP,
				HostIP:        daemonsetData.Containers[i].Ports[j].HostIP,
				HostPort:      daemonsetData.Containers[i].Ports[j].HostPort,
			}
		}
		containers[i].Ports = ports
		//组装每个容器的卷挂载组
		mounts := make([]corev1.VolumeMount, len(daemonsetData.Containers[i].MontVolume))
		for k, _ := range daemonsetData.Containers[i].MontVolume {
			mounts[k] = corev1.VolumeMount{
				Name:      daemonsetData.Containers[i].MontVolume[k].Name,
				ReadOnly:  daemonsetData.Containers[i].MontVolume[k].ReadOnly,
				MountPath: daemonsetData.Containers[i].MontVolume[k].MountPath,
				SubPath:   daemonsetData.Containers[i].MontVolume[k].SubPath,
			}
		}
		containers[i].VolumeMounts = mounts
		//组装环境变量
		envs := make([]corev1.EnvVar, len(daemonsetData.Containers[i].Envs))
		for l, _ := range daemonsetData.Containers[i].Envs {
			envs[l] = corev1.EnvVar{
				Name:      daemonsetData.Containers[i].Envs[l].Name,
				Value:     daemonsetData.Containers[i].Envs[l].Value,
				ValueFrom: nil,
			}
		}
		containers[i].Env = envs
		containers[i].ImagePullPolicy = corev1.PullPolicy(daemonsetData.Containers[i].ImagePullpolicy)
	}
	daemonset.Spec.Template.Spec.Containers = containers
	//判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
	if daemonsetData.HealthCheck {
		//设置容器的ReadinessProbe
		//若pod中有多个容器，则这里需要使用for循环去定义了
		for i, _ := range daemonset.Spec.Template.Spec.Containers {
			daemonset.Spec.Template.Spec.Containers[i].ReadinessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: daemonsetData.HealthPath,
						//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
						//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
						//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
						Port: intstr.IntOrString{
							Type:   0,
							IntVal: daemonsetData.Containers[i].Ports[i].ContainerPort,
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
			daemonset.Spec.Template.Spec.Containers[i].LivenessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: daemonsetData.HealthPath,
						Port: intstr.IntOrString{
							Type:   0,
							IntVal: daemonsetData.Containers[i].Ports[i].ContainerPort,
						},
					},
				},
				InitialDelaySeconds: 15,
				TimeoutSeconds:      5,
				PeriodSeconds:       5,
			}
		}
	}
	//当cpu和mem值不为空的时候，才去配置资源限制
	if daemonsetData.Mem != "" && daemonsetData.Cpu != "" {
		for i, _ := range daemonset.Spec.Template.Spec.Containers {
			//定义容器的limit和request资源: 设置 CPU 和内存的值
			daemonset.Spec.Template.Spec.Containers[i].Resources.Limits =

				map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(daemonsetData.Cpu),
					corev1.ResourceMemory: resource.MustParse(daemonsetData.Mem),
				}
			daemonset.Spec.Template.Spec.Containers[i].Resources.Requests =

				map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(daemonsetData.Cpu),
					corev1.ResourceMemory: resource.MustParse(daemonsetData.Mem),
				}
		}
	}
	fmt.Println("创建之前：", daemonset)
	//调用sdk创建deployment
	_, err = K8s.ClientSet.AppsV1().DaemonSets(daemonset.Namespace).Create(context.TODO(), daemonset, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建daemonset失败: " + err.Error())
		return errors.New("创建daemonset失败: " + err.Error())
	}
	return nil
}
