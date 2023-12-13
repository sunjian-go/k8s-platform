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

var StatefulSet statefulSet

type statefulSet struct {
}
type statefulSetResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

// 定义StatefulSet结构体，用于创建StatefulSet需要的参数属性的定义
type StatefulSetCreate struct {
	Name              string            `json:"name" binding:"required"`
	Namespace         string            `json:"namespace" binding:"required"`
	Replicas          int32             `json:"replicas"`
	Label             map[string]string `json:"label"`
	NodeName          string            `json:"nodeName"`
	NodeSelectorLabel map[string]string `json:"nodeSelectorLabel"`
	Containers        []Container       `json:"containers"`
	Volume            []Volume          `json:"volume"`
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

//	func (s *statefulSet) CreateStatefulSet(StatefulSetData *StatefulSetCreate) (err error) {
//		StatefulSet := &appsv1.StatefulSet{
//			ObjectMeta: metav1.ObjectMeta{
//				Name:      StatefulSetData.Name,
//				Namespace: StatefulSetData.Namespace,
//				Labels:    StatefulSetData.Label,
//			},
//			Spec: appsv1.StatefulSetSpec{
//				Replicas: &StatefulSetData.Replicas,
//				Selector: &metav1.LabelSelector{
//					MatchLabels: StatefulSetData.Label,
//				},
//				Template: corev1.PodTemplateSpec{
//					ObjectMeta: metav1.ObjectMeta{
//						Name:   StatefulSetData.Name,
//						Labels: StatefulSetData.Label,
//					},
//				},
//			},
//			Status: appsv1.StatefulSetStatus{},
//		}
//
//		//先准备一个container数组准备添加数据
//		containers := make([]corev1.Container, len(StatefulSetData.Containers))
//		for i, _ := range StatefulSetData.Containers {
//			//组装每个容器需要的端口配置
//			containers[i] = corev1.Container{
//				Name:  StatefulSetData.Containers[i].Name,
//				Image: StatefulSetData.Containers[i].Image,
//			}
//
//			//组装每个容器需要的端口组，查看一共几组ports，组装每一组，最后赋值给containers[i].Ports
//			ports := make([]corev1.ContainerPort, len(StatefulSetData.Containers[i].Ports))
//			for j, _ := range StatefulSetData.Containers[i].Ports {
//				ports[j] = corev1.ContainerPort{
//					Name:          StatefulSetData.Containers[i].Ports[j].PortName,
//					ContainerPort: StatefulSetData.Containers[i].Ports[j].ContainerPort,
//					Protocol:      corev1.ProtocolTCP,
//					HostIP:        StatefulSetData.Containers[i].Ports[j].HostIP,
//					HostPort:      StatefulSetData.Containers[i].Ports[j].HostPort,
//				}
//			}
//			containers[i].Ports = ports
//
//			//组装每个容器的卷挂载组
//			mounts := make([]corev1.VolumeMount, len(StatefulSetData.Containers[i].MontVolume))
//			for k, _ := range StatefulSetData.Containers[i].MontVolume {
//				mounts[k] = corev1.VolumeMount{
//					Name:      StatefulSetData.Containers[i].MontVolume[k].Name,
//					ReadOnly:  StatefulSetData.Containers[i].MontVolume[k].ReadOnly,
//					MountPath: StatefulSetData.Containers[i].MontVolume[k].MountPath,
//					SubPath:   StatefulSetData.Containers[i].MontVolume[k].SubPath,
//				}
//			}
//			containers[i].VolumeMounts = mounts
//
//			//组装环境变量
//			envs := make([]corev1.EnvVar, len(StatefulSetData.Containers[i].Envs))
//			for l, _ := range StatefulSetData.Containers[i].Envs {
//				envs[l] = corev1.EnvVar{
//					Name:      StatefulSetData.Containers[i].Envs[l].Name,
//					Value:     StatefulSetData.Containers[i].Envs[l].Value,
//					ValueFrom: nil,
//				}
//			}
//			containers[i].Env = envs
//			containers[i].ImagePullPolicy = corev1.PullPolicy(StatefulSetData.Containers[i].ImagePullpolicy)
//
//			//判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
//			if StatefulSetData.Containers[i].HealthCheck {
//				//设置容器的ReadinessProbe
//				//若pod中有多个容器，则这里需要使用for循环去定义了
//				for k, _ := range containers[i].Ports {
//					containers[i].ReadinessProbe = &corev1.Probe{
//						ProbeHandler: corev1.ProbeHandler{
//							HTTPGet: &corev1.HTTPGetAction{
//								Path: StatefulSetData.Containers[i].HealthPath,
//								//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
//								//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
//								//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
//								Port: intstr.IntOrString{
//									Type:   0,
//									IntVal: containers[i].Ports[k].ContainerPort,
//								},
//							},
//						},
//						//初始化等待时间
//						InitialDelaySeconds: 5,
//						//超时时间
//						TimeoutSeconds: 5,
//						//执行间隔
//						PeriodSeconds: 5,
//					}
//					containers[i].LivenessProbe = &corev1.Probe{
//						ProbeHandler: corev1.ProbeHandler{
//							HTTPGet: &corev1.HTTPGetAction{
//								Path: StatefulSetData.Containers[i].HealthPath,
//								//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
//								//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
//								//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
//								Port: intstr.IntOrString{
//									Type:   0,
//									IntVal: containers[i].Ports[k].ContainerPort,
//								},
//							},
//						},
//						InitialDelaySeconds: 15,
//						TimeoutSeconds:      5,
//						PeriodSeconds:       5,
//					}
//				}
//			}
//
//			//当cpu和mem值不为空的时候，才去配置资源限制
//			if StatefulSetData.Containers[i].Mem != "" && StatefulSetData.Containers[i].Cpu != "" {
//				//定义容器的limit和request资源: 设置 CPU 和内存的值
//				containers[i].Resources.Limits =
//					map[corev1.ResourceName]resource.Quantity{
//						corev1.ResourceCPU:    resource.MustParse(StatefulSetData.Containers[i].Cpu),
//						corev1.ResourceMemory: resource.MustParse(StatefulSetData.Containers[i].Mem),
//					}
//				containers[i].Resources.Reques ts =
//					map[corev1.ResourceName]resource.Quantity{
//						corev1.ResourceCPU:    resource.MustParse(StatefulSetData.Containers[i].Cpu),
//						corev1.ResourceMemory: resource.MustParse(StatefulSetData.Containers[i].Mem),
//					}
//			}
//		}
//		StatefulSet.Spec.Template.Spec.Containers = containers
//
//		//判断使用使用了挂载卷
//		if len(StatefulSetData.Volume) > 0 {
//			volumes := make([]corev1.Volume, len(StatefulSetData.Volume))
//			var volumeSource corev1.VolumeSource
//			for i, _ := range StatefulSetData.Volume {
//				//遍历单个容器中有几组卷，根据卷类型，挨个组装完整的单个卷对象
//				switch StatefulSetData.Volume[i].Type {
//				case "configMap":
//					volumeSource = corev1.VolumeSource{
//						ConfigMap: &corev1.ConfigMapVolumeSource{
//							LocalObjectReference: corev1.LocalObjectReference{
//								Name: StatefulSetData.Volume[i].Context,
//							},
//							DefaultMode: &StatefulSetData.Volume[i].Mode,
//						},
//					}
//				case "HostPath":
//					volumeSource = corev1.VolumeSource{
//						HostPath: &corev1.HostPathVolumeSource{
//							Path: StatefulSetData.Volume[i].Context,
//						},
//					}
//				case "EmptyDir":
//					volumeSource = corev1.VolumeSource{
//						EmptyDir: &corev1.EmptyDirVolumeSource{},
//					}
//				case "PersistentVolumeClaim":
//					volumeSource = corev1.VolumeSource{
//						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
//							ClaimName: StatefulSetData.Volume[i].Context,
//						},
//					}
//				case "Secret":
//					volumeSource = corev1.VolumeSource{
//						Secret: &corev1.SecretVolumeSource{
//							SecretName: StatefulSetData.Volume[i].Context,
//						},
//					}
//				}
//				//给volume数组赋值
//				volumes[i] = corev1.Volume{
//					Name:         StatefulSetData.Volume[i].VolumeName,
//					VolumeSource: volumeSource,
//				}
//				fmt.Println("赋值前的volume为：", volumes[i])
//			}
//			StatefulSet.Spec.Template.Spec.Volumes = volumes
//			fmt.Println("卷数据为：", StatefulSet.Spec.Template.Spec.Volumes)
//		}
//
//		//判断是否使用节点亲和性和nodeName
//		if StatefulSetData.NodeSelectorLabel != nil {
//			StatefulSet.Spec.Template.Spec.NodeSelector = StatefulSetData.NodeSelectorLabel
//		}
//		if StatefulSetData.NodeName != "" {
//			StatefulSet.Spec.Template.Spec.NodeName = StatefulSetData.NodeName
//		}
//		fmt.Println("创建之前：", StatefulSet.Spec.Template.Spec.Volumes[0].Name, StatefulSet.Spec.Template.Spec.Volumes[1].Name)
//		//调用sdk创建deployment
//		_, err = K8s.ClientSet.AppsV1().StatefulSets(StatefulSet.Namespace).Create(context.TODO(), StatefulSet, metav1.CreateOptions{})
//		if err != nil {
//			logger.Error("创建StatefulSet失败: " + err.Error())
//			return errors.New("创建StatefulSet失败: " + err.Error())
//		}
//		return nil
//	}
func (s *statefulSet) CreateStatefulSet(StatefulSetData *StatefulSetCreate) (err error) {
	fmt.Println("传入的数据为：", StatefulSetData)
	labels := make(map[string]string)
	//先添加默认标签
	labels["Controller"] = StatefulSetData.Name + "-" + StatefulSetData.Namespace
	//如果有自定义标签，就加上自定义标签
	if len(StatefulSetData.Label) > 0 {
		for key, value := range StatefulSetData.Label {
			labels[key] = value
		}
	}
	StatefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      StatefulSetData.Name,
			Namespace: StatefulSetData.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   StatefulSetData.Name,
					Labels: labels,
				},
			},
			Replicas: &StatefulSetData.Replicas,
		},
		Status: appsv1.StatefulSetStatus{},
	}

	//判断使用使用了挂载卷
	if len(StatefulSetData.Volume) > 0 {
		volumes := make([]corev1.Volume, len(StatefulSetData.Volume))
		var volumeSource corev1.VolumeSource
		for i, _ := range StatefulSetData.Volume {
			//遍历单个容器中有几组卷，根据卷类型，挨个组装完整的单个卷对象
			switch StatefulSetData.Volume[i].Type {
			case "configMap":
				//如果Items不为空就说明有特定键，就去赋值特定键即可
				if StatefulSetData.Volume[i].Items != nil {
					items := make([]corev1.KeyToPath, len(StatefulSetData.Volume[i].Items))
					for k, _ := range StatefulSetData.Volume[i].Items {
						items[k].Key = StatefulSetData.Volume[i].Items[k].Key
						items[k].Path = StatefulSetData.Volume[i].Items[k].Path
						items[k].Mode = &StatefulSetData.Volume[i].Items[k].Mode
					}
					volumeSource = corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: StatefulSetData.Volume[i].Context,
							},
							Items:       items,
							DefaultMode: &StatefulSetData.Volume[i].Mode,
						},
					}
				} else { //否则就正常赋值
					volumeSource = corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: StatefulSetData.Volume[i].Context,
							},
							DefaultMode: &StatefulSetData.Volume[i].Mode,
						},
					}
				}
			case "HostPath":
				hostPathtype := corev1.HostPathType(StatefulSetData.Volume[i].HostPathType)
				volumeSource = corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: StatefulSetData.Volume[i].Context,
						Type: &hostPathtype,
					},
				}
			case "EmptyDir":
				volumeSource = corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				}
			case "PersistentVolumeClaim":
				volumeSource = corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: StatefulSetData.Volume[i].Context,
					},
				}
			case "Secret":
				volumeSource = corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: StatefulSetData.Volume[i].Context,
					},
				}
			}
			//给volume数组赋值
			volumes[i] = corev1.Volume{
				Name:         StatefulSetData.Volume[i].VolumeName,
				VolumeSource: volumeSource,
			}
			fmt.Println("赋值前的volume为：", volumes[i])
		}
		StatefulSet.Spec.Template.Spec.Volumes = volumes
		fmt.Println("卷数据为：", StatefulSet.Spec.Template.Spec.Volumes)
	}

	//先准备一个container数组准备添加数据
	containers := make([]corev1.Container, len(StatefulSetData.Containers))
	for i, _ := range StatefulSetData.Containers {
		//基本配置组组装
		containers[i] = corev1.Container{
			Name:  StatefulSetData.Containers[i].Name,
			Image: StatefulSetData.Containers[i].Image,
		}

		//组装每个容器需要的端口组，查看一共几组ports，组装每一组，最后赋值给containers[i].Ports
		if len(StatefulSetData.Containers[i].Ports) > 0 {
			ports := make([]corev1.ContainerPort, len(StatefulSetData.Containers[i].Ports))
			for j, _ := range StatefulSetData.Containers[i].Ports {
				ports[j] = corev1.ContainerPort{
					Name:          StatefulSetData.Containers[i].Ports[j].PortName,
					ContainerPort: StatefulSetData.Containers[i].Ports[j].ContainerPort,
					//HostIP:        StatefulSetData.Containers[i].Ports[j].HostIP,
					//HostPort:      StatefulSetData.Containers[i].Ports[j].HostPort,
				}
				if StatefulSetData.Containers[i].Ports[j].Protocol != "" {
					switch StatefulSetData.Containers[i].Ports[j].Protocol {
					case "TCP":
						ports[j].Protocol = corev1.ProtocolTCP
					case "UDP":
						ports[j].Protocol = corev1.ProtocolUDP

					}
				} else {
					ports[j].Protocol = corev1.ProtocolTCP
				}
				if StatefulSetData.Containers[i].Ports[j].HostportStatus {
					ports[j].HostIP = StatefulSetData.Containers[i].Ports[j].HostIP
					ports[j].HostPort = StatefulSetData.Containers[i].Ports[j].HostPort
				}
			}
			containers[i].Ports = ports
		}

		//组装每个容器的卷挂载组
		if len(StatefulSetData.Containers[i].MontVolume) > 0 {
			mounts := make([]corev1.VolumeMount, len(StatefulSetData.Containers[i].MontVolume))
			for k, _ := range StatefulSetData.Containers[i].MontVolume {
				mounts[k] = corev1.VolumeMount{
					Name:      StatefulSetData.Containers[i].MontVolume[k].Name,
					ReadOnly:  StatefulSetData.Containers[i].MontVolume[k].ReadOnly,
					MountPath: StatefulSetData.Containers[i].MontVolume[k].MountPath,
					SubPath:   StatefulSetData.Containers[i].MontVolume[k].SubPath,
				}
			}
			containers[i].VolumeMounts = mounts
		}

		//组装环境变量
		if len(StatefulSetData.Containers[i].Envs) > 0 {
			envs := make([]corev1.EnvVar, len(StatefulSetData.Containers[i].Envs))
			for l, _ := range StatefulSetData.Containers[i].Envs {
				envs[l] = corev1.EnvVar{
					Name:      StatefulSetData.Containers[i].Envs[l].Name,
					Value:     StatefulSetData.Containers[i].Envs[l].Value,
					ValueFrom: nil,
				}
			}
			containers[i].Env = envs
		}
		containers[i].ImagePullPolicy = corev1.PullPolicy(StatefulSetData.Containers[i].ImagePullpolicy)

		//判断是否打开健康检查功能，若打开，则定义ReadinessProbe和LivenessProbe
		if StatefulSetData.Containers[i].HealthCheck {
			//设置容器的ReadinessProbe
			//若pod中有多个容器，则这里需要使用for循环去定义了
			//for k, _ := range containers[i].Ports {
			containers[i].ReadinessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: StatefulSetData.Containers[i].HealthPath,
						//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
						//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
						//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
						Port: intstr.IntOrString{
							Type:   0,
							IntVal: int32(80),
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
			containers[i].LivenessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: StatefulSetData.Containers[i].HealthPath,
						//intstr.IntOrString的作用是端口可以定义为整型，也可以定义为字符串
						//Type=0则表示表示该结构体实例内的数据为整型，转json时只使用IntVal的数据
						//Type=1则表示表示该结构体实例内的数据为字符串，转json时只使用StrVal的数据
						Port: intstr.IntOrString{
							Type:   0,
							IntVal: int32(80),
						},
					},
				},
				InitialDelaySeconds: 15,
				TimeoutSeconds:      5,
				PeriodSeconds:       5,
			}
			//}
		}

		//当请求资源cpu和mem值不为空的时候，才去配置资源
		if StatefulSetData.Containers[i].ReqMem != "" && StatefulSetData.Containers[i].ReqCpu != "" {
			containers[i].Resources.Requests =
				map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(StatefulSetData.Containers[i].ReqCpu),
					corev1.ResourceMemory: resource.MustParse(StatefulSetData.Containers[i].ReqMem),
				}
		} else {
			if StatefulSetData.Containers[i].ReqMem != "" {
				containers[i].Resources.Requests =
					map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceMemory: resource.MustParse(StatefulSetData.Containers[i].ReqMem),
					}
			}
			if StatefulSetData.Containers[i].ReqCpu != "" {
				containers[i].Resources.Requests =
					map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceCPU: resource.MustParse(StatefulSetData.Containers[i].ReqCpu),
					}
			}
		}

		//当限制cpu和mem值不为空的时候，才去配置资源限制
		if StatefulSetData.Containers[i].Mem != "" && StatefulSetData.Containers[i].Cpu != "" {
			//定义容器的limit和request资源: 设置 CPU 和内存的值
			containers[i].Resources.Limits =
				map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse(StatefulSetData.Containers[i].Cpu),
					corev1.ResourceMemory: resource.MustParse(StatefulSetData.Containers[i].Mem),
				}
		} else {
			if StatefulSetData.Containers[i].Mem != "" {
				containers[i].Resources.Limits =
					map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceMemory: resource.MustParse(StatefulSetData.Containers[i].Mem),
					}
			}
			if StatefulSetData.Containers[i].Cpu != "" {
				containers[i].Resources.Limits =
					map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceCPU: resource.MustParse(StatefulSetData.Containers[i].Cpu),
					}
			}
		}
	}
	StatefulSet.Spec.Template.Spec.Containers = containers

	//判断是否使用节点亲和性
	if StatefulSetData.NodeSelectorLabel != nil {
		StatefulSet.Spec.Template.Spec.NodeSelector = StatefulSetData.NodeSelectorLabel
	}
	if StatefulSetData.NodeName != "" {
		StatefulSet.Spec.Template.Spec.NodeName = StatefulSetData.NodeName
	}
	//fmt.Println("创建之前：", StatefulSet.Spec.Template.Spec.Volumes[3].Name)
	//调用sdk创建deployment
	_, err = K8s.ClientSet.AppsV1().StatefulSets(StatefulSet.Namespace).Create(context.TODO(), StatefulSet, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建StatefulSet失败: " + err.Error())
		return errors.New("创建StatefulSet失败: " + err.Error())
	}
	return nil
}
