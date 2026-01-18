package main

import (
	"goStudyHm4/config"
	"goStudyHm4/model"
	"goStudyHm4/router"
	"goStudyHm4/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化日志
	utils.InitLogger()
	logrus.Info("博客系统后端启动...")

	// 初始化数据库
	config.InitDB()

	// 数据库表迁移
	model.AutoMigrate()

	// 初始化路由
	r := router.InitRouter()

	// 启动服务
	logrus.Info("服务启动成功，监听端口: 8080")
	if err := r.Run(":8080"); err != nil {
		logrus.Fatalf("服务启动失败: %v", err)
	}
}
