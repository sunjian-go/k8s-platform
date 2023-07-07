package controller

import "github.com/gin-gonic/gin"

var Router router

type router struct {
}

func (r *router) InitApiRouter(router *gin.Engine) {

	router.
		//pod路由
		GET("/api/corev1/pods", Pod.GetPods).
		GET("/api/corev1/podetail", Pod.GetPodDetail).
		DELETE("/api/corev1/deletepod", Pod.DeletePod).
		PUT("/api/corev1/updatepod", Pod.UpdatePod).
		GET("/api/corev1/getcontainers", Pod.GetContainer).
		GET("/api/corev1/getlog", Pod.GetContainerLog).
		GET("/api/corev1/getpodnum", Pod.GetNamespacePod).
		//deployment路由
		GET("/api/appsv1/getdeployments", Deployment.GetDeployments).
		GET("/api/appsv1/detdeploymentdetail", Deployment.GetDeploymentDetail).
		DELETE("/api/appsv1/deletedeployment", Deployment.DeleteDeployment).
		PUT("/api/appsv1/scaledeployment", Deployment.ScaleDeployment).
		POST("/api/appsv1/createdeployment", Deployment.CreateDeployment).
		PUT("/api/appsv1/restartdeployment", Deployment.RestartDeployment).
		PUT("/api/appsv1/updatedeployment", Deployment.UpdateDeployment).
		GET("/api/appsv1/getnamespacedeployNum", Deployment.GetNamespaceDeployNum).
		//daemonSet路由
		GET("/api/appsv1/getdaemonSet", DaemonSet.GetDaemonSets).
		GET("/api/appsv1/getdaemonSetDetail", DaemonSet.GetDaemonSetDetail).
		DELETE("/api/appsv1/deleteDaemonSet", DaemonSet.DeleteDaemonSet).
		PUT("/api/appsv1/updateDaemonSet", DaemonSet.UpdateDaemonSet).
		//StatefulSet路由
		GET("/api/appsv1/getstatefulSets", StatefulSet.GetStatefulSets).
		GET("/api/appsv1/getstatefulSetDetail", StatefulSet.GetStatefulSetDetail).
		DELETE("/api/appsv1/deletestatefulSet", StatefulSet.DeleteStatefulSet).
		PUT("/api/appsv1/updatestatefulSet", StatefulSet.UpdateStatefulSet)

}
