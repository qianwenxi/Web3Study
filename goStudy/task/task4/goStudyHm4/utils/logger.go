package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// 初始化日志配置
func InitLogger() {
	// 设置日志格式为JSON格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 设置日志输出到控制台和文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Warn("无法打开日志文件，将日志输出到控制台")
	}
	// 设置日志级别（Debug级别，输出所有日志）
	logrus.SetLevel(logrus.DebugLevel)
}
