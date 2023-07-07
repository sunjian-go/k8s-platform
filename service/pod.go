package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	"k8s-platform/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 定义pod类型和Pod对象，用于包外的调用(包是指service目录)，例如Controller
var Pod pod

type pod struct {
}

// 定义列表的返回内容，Items是pod元素列表，Total为pod元素数量
type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}
type PodsNp struct {
	Namespace string
	PodNum    int
}

// 定义PodsNp类型，用于返回namespace中pod的数量
type PodNp struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"podNum"`
}

// toCells方法用于将pod类型数组，转换成DataCell类型数组
func (p *pod) toCell(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

// fromCells方法用于将DataCell类型数组，转换成pod类型数组
func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		//cells[i].(podCell)就使用到了断言,断言后转换成了podCell类型，然后又转换成了Pod类型
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

// 获取pod列表
// 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsresp *PodsResp, err error) {
	//fmt.Println(filterName, namespace, limit, page)
	//获取podList类型的pod列表
	//context.TODO()用于声明一个空的context上下文，用于List方法内设置这个请求的超时（源码），这里的常用用法
	//metav1.ListOptions{}用于过滤List数据，如使用label，field等
	//kubectl get services --all-namespaces --field-seletor metadata.namespace != default
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	//fmt.Println("podlist: ", podList.Items)
	if err != nil {
		//logger用于打印日志
		//return用于返回response内容
		logger.Info("获取pod列表失败" + err.Error())
		return nil, errors.New("获取pod列表失败" + err.Error())
	}
	//实例化dataSelector对象
	selecttableData := &DataSelector{
		GenericDataList: p.toCell(podList.Items), //将pods列表转换为DataCell类型赋值
		DataSelectQuery: &DataSelectQuery{ //
			FilterQuery: &FilterQuery{
				Name: filterName, //将传进来的需要查找的Name赋值给该结构体
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit, //将传进来的页数和每页的数量赋值
				Page:  page,
			},
		},
	}
	//先过滤
	filterd := selecttableData.Filter()
	//total := len(selecttableData.GenericDataList)
	total := len(filterd.GenericDataList) //计算过滤好的目标pod列表的长度
	//fmt.Println("过滤后：", filterd.GenericDataList)
	//再排序和分页
	//for _, data := range filterd.GenericDataList {
	//	fmt.Println("排序分页前：", data.GetName(), data.GetCreation())
	//}
	pods := filterd.Sort().Paginate() //连续调用排序和分页方法
	//for _, data := range pods.GenericDataList {
	//	fmt.Println("排序分页后：", data.GetName(), data.GetCreation())
	//}
	//将[]DataCell类型的pod列表转为v1.pod列表
	podv1s := p.fromCells(pods.GenericDataList)
	//fmt.Println("整理好的pod信息：", podv1s)
	return &PodsResp{
		Items: podv1s,
		Total: total,
	}, nil
}

// 获取pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Info("获取pod详情失败: ", err.Error())
		return nil, errors.New("获取pod详情失败: " + err.Error())
	}
	return pod, nil
}

// 删除pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	//K8s.ClientSet.CoreV1().Pods(namespace).Create(context.TODO(), podName, metav1.CreateOptions{})
	if err != nil {
		logger.Error("删除pod失败：", err.Error())
		return errors.New("删除pod失败：" + err.Error())
	}
	return nil
}

// 更新pod
func (p *pod) UpdatePod(podName, namespace, content string) (err error) {
	pod := &corev1.Pod{}
	if err = json.Unmarshal([]byte(content), pod); err != nil {
		logger.Error("反序列化失败: ", err.Error())
		return errors.New("反序列化失败: " + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新pod"+podName+"失败: ", err.Error())
		return errors.New("更新pod" + podName + "失败: " + err.Error())
	}
	return nil
}

// 获取容器信息
func (p *pod) GetContainer(podName, namespace string) (containers []string, err error) {
	//获取pod详情
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		logger.Info("获取pod详情失败: ", err.Error())
		return nil, errors.New("获取pod详情失败: " + err.Error())
	}
	//从pod中拿到容器名
	for _, cont := range pod.Spec.Containers {
		containers = append(containers, cont.Name)
	}
	return containers, nil
}

// 获取pod内容器日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	//1.设置日志的配置，容器名，获取的内容的配置
	lineLimit := int64(config.PodLogTailLine) //先将定义的行数转为int64位
	option := &corev1.PodLogOptions{          //定义一个corev1.PodLogOptions指针并赋值
		Container: containerName,
		TailLines: &lineLimit,
	}
	//2.获取一个request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	//3.发起stream连接，得到Response.body
	podLog, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取podLog失败" + err.Error()))
		return "", errors.New("获取podLog失败" + err.Error())
	}
	defer podLog.Close() //关闭stream连接
	//4.将response.body写入到缓冲区，目的是为了转换成string类型
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLog) //将得到的Log数据存在缓冲区
	if err != nil {
		logger.Error(errors.New("拷贝podLog失败" + err.Error()))
		return "", errors.New("拷贝podLog失败" + err.Error())
	}
	//5.转换数据返回
	return buf.String(), nil
}

// 获取每个namespace中pod的数量
func (p *pod) GetNamespacePod() (podsNps []*PodsNp, err error) {
	//获取namespace列表
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("获取namespace信息失败：" + err.Error())
		return nil, errors.New("获取namespace信息失败：" + err.Error())
	}
	//获取pod列表
	for _, namespace := range namespaceList.Items {
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error("获取pod信息失败：" + err.Error())
			return nil, errors.New("获取pod信息失败：" + err.Error())
		}
		//组装数据
		podsnp := &PodsNp{
			PodNum:    len(podList.Items),
			Namespace: namespace.Name,
		}
		//添加到podsNps数组中
		podsNps = append(podsNps, podsnp)
	}
	return podsNps, nil
}
