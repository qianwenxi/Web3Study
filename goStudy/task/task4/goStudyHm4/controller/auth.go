package controller

import (
	"net/http"

	"goStudyHm4/config"
	"goStudyHm4/model"
	"goStudyHm4/utils"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("注册参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		logrus.Errorf("创建用户失败: %v", err)
		c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	logrus.Infof("用户注册成功: %s", req.Username)
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功", "user_id": user.ID})
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("登录参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 查询用户
	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		logrus.Warnf("用户不存在: %s", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.Warnf("密码错误: %s", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		logrus.Errorf("生成Token失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	logrus.Infof("用户登录成功: %s", req.Username)
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
