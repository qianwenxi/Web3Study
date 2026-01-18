package middleware

import (
	"net/http"

	"goStudyHm4/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthRequired 需认证的接口中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logrus.Warn("请求缺少Authorization头")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			c.Abort()
			return
		}

		// 解析Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			logrus.Warnf("Token解析失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token或Token已过期"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
