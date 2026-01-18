package config

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// JWT配置
var JWTSecret = []byte("your_secure_secret_key_123") // 生产环境请更换为随机密钥
var TokenExpireHour = 24                             // Token有效期（小时）

// 数据库实例
var DB *gorm.DB

// 初始化数据库连接
func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/WEB3STUDY?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		logrus.Fatalf("数据库连接失败: %v", err)
	}
	logrus.Info("数据库连接成功")
}
