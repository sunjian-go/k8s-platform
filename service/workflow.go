package service

import (
	"k8s-platform/dao"
	"k8s-platform/model"
)

var Workflow workflow

type workflow struct {
}

// 定义WorkflowCreate结构体，用于创建workflow需要的参数属性的定义
type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int32                  `json:"replicas"`
	Image         string                 `json:"image"`
	Label         map[string]string      `json:"label"`
	Cpu           string                 `json:"cpu"`
	Memory        string                 `json:"memory"`
	ContainerPort int32                  `json:"container_port"`
	HealthCheck   bool                   `json:"health_check"`
	HealthPath    string                 `json:"health_path"`
	Type          string                 `json:"type"`
	Port          int32                  `json:"port"`
	NodePort      int32                  `json:"node_port"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
}

// 获取列表分页查询
func (w *workflow) GetList(name string, page, limit int) (data *dao.WorkflowResp, err error) {
	data, err = dao.Workflow.GetWorkflow(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 查询workflow单条数据
func (w *workflow) GetById(id int) (data *model.WorkFlow, err error) {
	data, err = dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 创建workflow
// workflow的类型分为三种，ClusterIP、NodePort、Ingress
func (w *workflow) CreateWorkflow(data *WorkflowCreate) (err error) {
	//若workflow不是ingress类型，传入空字符串即可
	var ingressName string
	//为了判断是否需要新增ingress
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}
	//组装mysql中workflow的单条数据
	workflow := &model.WorkFlow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}
	//调用dao层执行数据库的添加操作
	err = dao.Workflow.Add(workflow)
	if err != nil {
		return err
	}
	//创建k8s资源
	err = createWorkflowResp(data)
	if err != nil {
		return err
	}
	return nil

}
func createWorkflowResp(data *WorkflowCreate) (err error) {
	//创建deployment
	dc := &DeployCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Img:           data.Image,
		Replicas:      data.Replicas,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Mem:           data.Memory,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}
	err = Deployment.CreateDeployment(dc)
	if err != nil {
		return err
	}

	var svcType string
	if data.Type != "Ingress" {
		svcType = data.Type //如果类型不是Ingress，那么类型是啥就是啥
	} else {
		svcType = "ClusterIP" //如果类型是Ingress，就将svc的类型改为ClusterIP,因为ingress不需要svc暴露端口
	}
	//创建svc
	sc := &ServiceCreate{
		Name:          getServiceName(data.Name),
		Namespace:     data.Namespace,
		Type:          svcType,
		ContainerPort: data.ContainerPort,
		Port:          data.Port,
		NodePort:      data.NodePort,
		Label:         data.Label,
	}
	err = SVC.CreateSvc(sc)
	if err != nil {
		return err
	}
	//创建ing
	if data.Type == "Ingress" {
		ic := &IngressCreate{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		err = Ingress.CreateIngress(ic)
		if err != nil {
			return err
		}
	}

	return nil
}

// workflow名字转换成service名字，添加-svc后缀
func getServiceName(workflowName string) (serviceName string) {
	return workflowName + "-svc"
}

// workflow名字转换成ingress名字，添加-ing后缀
func getIngressName(workflowName string) (ingressName string) {
	return workflowName + "-ing"
}

// 删除workflow
func (w *workflow) DelById(id int) (err error) {
	//获取数据库数据,用于删除k8s资源的参数
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}
	//删除k8s资源
	err = delWorkflowResp(workflow)
	if err != nil {
		return err
	}
	//删除数据库数据
	err = dao.Workflow.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// 删除k8s资源 deployment service ingress
func delWorkflowResp(data *model.WorkFlow) (err error) {
	err = Ingress.DeleteIngress(getIngressName(data.Name), data.Namespace)
	if err != nil {
		return err
	}
	err = SVC.DeleteSvc(getServiceName(data.Name), data.Namespace)
	if err != nil {
		return err
	}
	err = Deployment.DeleteDeployment(data.Name, data.Namespace)
	if err != nil {
		return err
	}
	return nil
}
