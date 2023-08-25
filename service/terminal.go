package service

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct {
}

// ws处理逻辑函数
func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	//加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("加载k8s配置失败：" + err.Error())
		return
	}
	//解析form入参，获取namespace,pod,container参数
	if err := r.ParseForm(); err != nil {
		logger.Error("解析参数失败：" + err.Error())
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	logger.Info("exec pod: %s, container: %s, namespace: %s\n", podName, containerName, namespace)

	//new一个TerminalSession类型的pty实例
	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		logger.Error("get pty failed: %v\n", err)
		return
	}
	//处理关闭
	defer func() {
		logger.Info("close session.")
		pty.Close()
	}()

	//组装post请求
	// 初始化pod所在的corev1资源组
	// PodExecOptions struct 包括Container stdout stdout Command 等结构
	// scheme.ParameterCodec 应该是pod 的GVK （GroupVersion & Kind）之类的
	// URL长相:
	// https://192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c7rmsk/exec?
	//command=%2Fbin%2Fbash&container=nginxwf2&stderr=true&stdin=true&stdout=true&tty=true
	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	logger.Info("exec post request url: ", req)

	//升级SPDY协议
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		logger.Error("建立SPDY连接失败，" + err.Error())
		return
	}
	//与kublet建立stream连接
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
	if err != nil {
		logger.Error("执行pod命令失败，" + err.Error())
		//将报错返回出去
		pty.Write([]byte("执行pod命令失败，" + err.Error()))
		//标记退出stream流
		pty.Done()
	}
}

// 消息内容
type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      string `json:"rows"`
	Cols      string `json:"cols"`
}

// 交互的结构体，接管输入和输出
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// 初始化一个websockerf.Upgrader类型的对象，用于http协议升级为ws协议
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// 创建TerminalSession类型的对象并返回
func NewTerminalSession(w http.ResponseWriter, r http.Request, responseHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, errors.New("升级websocket失败：" + err.Error())
	}
	//new
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return session, nil
}

// 读消息
// 用于读取web端的输入，接收web端输入的指令内容,返回值int是读成功了多少数据
func (t *TerminalSession) Read(p []byte) (int, error) {
	//从ws中读取消息
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("read message err: %v", err)
		return 0, err
	}
	//从ws中读取出来的消息进行反序列化
	var msg terminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("read parse message err: %v", err)
		return 0, err
	}
	//根据消息内容的选项做不同动作
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("unknown message type '%s'", msg.Operation)
		return 0, fmt.Errorf("unknown message type '%s'", msg.Operation)
	}
}

// 写数据的方法
// 拿到apiserver的返回内容，向web端输出
func (t *TerminalSession) Write(p []byte) (int, error) {
	//将apiserver的返回内容组装进结构体并进行编码
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("组装消息结构体失败：%v", err)
		return 0, err
	}
	//开始写数据
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("写入消息失败：%v", err)
		return 0, err
	}
	//返回写入数据的长度
	return len(p), nil
}

// 标记关闭的方法
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// 关闭的方法
func (t *TerminalSession) Close() {
	t.wsConn.Close()
}

// resize方法，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan: //读取到size数据的话就返回该数据
		return &size
	case <-t.doneChan: //读取到数据的话就代表关闭ws
		return nil
	}
}
