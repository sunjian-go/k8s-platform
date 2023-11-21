package controller

import (
	"github.com/gin-gonic/gin"
)

var Router router

type router struct {
}

func (r *router) InitApiRouter(router *gin.Engine) {

	router.
		//pod路由
		POST("/api/login", Login.Login).
		GET("/api/corev1/getpods", Pod.GetPods).
		GET("/api/corev1/podetail", Pod.GetPodDetail).
		DELETE("/api/corev1/deletepod", Pod.DeletePod).
		PUT("/api/corev1/updatepod", Pod.UpdatePod).
		GET("/api/corev1/getcontainers", Pod.GetContainer).
		GET("/api/corev1/getlog", Pod.GetContainerLog).
		GET("/api/corev1/getpodnum", Pod.GetNamespacePod).
		//deployment路由
		GET("/api/appsv1/getdeployments", Deployment.GetDeployments).
		GET("/api/appsv1/getdeploymentdetail", Deployment.GetDeploymentDetail).
		DELETE("/api/appsv1/deletedeployment", Deployment.DeleteDeployment).
		PUT("/api/appsv1/scaledeployment", Deployment.ScaleDeployment).
		POST("/api/appsv1/createdeployment", Deployment.CreateDeployment).
		PUT("/api/appsv1/restartdeployment", Deployment.RestartDeployment).
		PUT("/api/appsv1/updatedeployment", Deployment.UpdateDeployment).
		GET("/api/appsv1/getnamespacedeployNum", Deployment.GetNamespaceDeployNum).
		//daemonSet路由
		GET("/api/appsv1/getdaemonSet", DaemonSet.GetDaemonSets).
		GET("/api/appsv1/getdaemonSetDetail", DaemonSet.GetDaemonSetDetail).
		POST("/api/appsv1/createDaemonset", DaemonSet.CreateDaemonSet).
		DELETE("/api/appsv1/deleteDaemonSet", DaemonSet.DeleteDaemonSet).
		PUT("/api/appsv1/updateDaemonSet", DaemonSet.UpdateDaemonSet).
		//StatefulSet路由
		GET("/api/appsv1/getstatefulSets", StatefulSet.GetStatefulSets).
		GET("/api/appsv1/getstatefulSetDetail", StatefulSet.GetStatefulSetDetail).
		POST("/api/appsv1/createStatefulSet", StatefulSet.CreateStatefulSet).
		DELETE("/api/appsv1/deletestatefulSet", StatefulSet.DeleteStatefulSet).
		PUT("/api/appsv1/updatestatefulSet", StatefulSet.UpdateStatefulSet).
		//node路由
		GET("/api/corev1/getnodes", Node.GetNodes).
		GET("/api/corev1/getnodedetail", Node.GetNodeDetail).
		//namespace路由
		GET("/api/corev1/getnamespaces", Namespace.GetNamespaces).
		GET("/api/corev1/getnamespaceList", Namespace.GetnamespaceList).
		GET("/api/corev1/getnamespaceDetail", Namespace.GetNamespaceDetail).
		DELETE("/api/corev1/deletenamespace", Namespace.DeleteNamespace).
		//PV路由
		GET("/api/corev1/getpvs", Pv.GetPvs).
		GET("/api/corev1/getpvdetail", Pv.GetPvDetail).
		DELETE("/api/corev1/deletepv", Pv.DeletePv).
		//svc路由
		GET("/api/corev1/getsvc", SVC.GetSvcs).
		GET("/api/corev1/getsvcdetail", SVC.GetSvcDetail).
		POST("/api/corev1/createsvc", SVC.CreateSvc).
		DELETE("/api/corev1/deletesvc", SVC.DeleteSvc).
		PUT("/api/corev1/updatesvc", SVC.UpdateSvc).
		//ingress路由
		GET("/api/networking/geting", Ingress.GetIngresses).
		GET("/api/networking/getingdetail", Ingress.GetIngressDetail).
		POST("/api/networking/createing", Ingress.CreateIngress).
		DELETE("/api/networking/deleteing", Ingress.DeleteIngress).
		PUT("/api/networking/updateing", Ingress.UpdateIngress).
		//configMap路由
		GET("/api/corev1/getcms", ConfigMap.GetConfigMaps).
		GET("/api/corev1/getcmdetail", ConfigMap.GetConfigMapDetail).
		DELETE("/api/corev1/deletecm", ConfigMap.DeleteConfigMap).
		PUT("/api/corev1/updatecm", ConfigMap.UpdateConfigMap).
		//secret路由
		GET("/api/corev1/getsecrets", Secret.GetSecrets).
		GET("/api/corev1/getsecretdetail", Secret.GetSecretDetail).
		POST("/api/corev1/createSecret", Secret.CreateSecret).
		DELETE("/api/corev1/deletesecret", Secret.DeleteSecret).
		PUT("/api/corev1/updatesecret", Secret.UpdateSecret).
		//PVC路由
		GET("/api/corev1/getpvcs", Pvc.GetPvcs).
		GET("/api/corev1/getpvcdetail", Pvc.GetPvcDetail).
		POST("/api/corev1/createPvc", Pvc.CreatePvc).
		DELETE("/api/corev1/deletepvc", Pvc.DeletePvc).
		PUT("/api/corev1/updatepvc", Pvc.UpdatePvc).
		//workflow路由
		GET("/api/workflow/getworkflows", Workflow.GetWorkflows).
		GET("/api/workflow/getbyid", Workflow.GetById).
		DELETE("/api/workflow/delbyid", Workflow.DelById).
		POST("/api/workflow/createworkflow", Workflow.CreateWorkflow).
		//resource路由
		GET("/api/resources", Resources.GetAllResources).
		//ws获取日志
		GET("/api/getlogs", Pod.GetLog).
		//前端设置背景色
		GET("/api/getColor", Styles.GetColor).
		PUT("/api/updateColor", Styles.UpdateColor)

}
