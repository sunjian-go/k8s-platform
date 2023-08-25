package config

import "time"

const (
	ListenAddr     = "0.0.0.0:8999"
	Kubeconfig     = "D:\\sunjian\\golang\\ncode\\k8s-platform\\k8sconf\\config.txt"
	PodLogTailLine = 2000 //tail的日志行数 tail -n 2000

	//数据库配置
	DbType = "mysql"
	DbHost = "127.0.0.1"
	DbPort = 3306
	DbName = "k8s_demo"
	DbUser = "root"
	DbPwd  = "Tsit@2022"

	//打印mysql debug的sql日志
	LogMode = true

	//连接池的配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
)
