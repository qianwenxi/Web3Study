package controller

import (
	"net/http"

	"goStudyHm4/config"
	"goStudyHm4/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreatePostRequest 创建文章请求参数
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=10"`
}

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("创建文章参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userID")

	// 创建文章
	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}
	if err := config.DB.Create(&post).Error; err != nil {
		logrus.Errorf("创建文章失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	logrus.Infof("用户%d创建文章成功: %s", userID.(uint), req.Title)
	c.JSON(http.StatusCreated, gin.H{"message": "文章创建成功", "post": post})
}

// GetPosts 获取所有文章列表
func GetPosts(c *gin.Context) {
	var posts []model.Post
	// 关联查询用户信息
	if err := config.DB.Preload("User").Find(&posts).Error; err != nil {
		logrus.Errorf("查询文章列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	logrus.Info("查询文章列表成功")
	c.JSON(http.StatusOK, gin.H{"posts": posts, "count": len(posts)})
}

// GetPost 获取单篇文章详情
func GetPost(c *gin.Context) {
	// 获取路径参数id
	postID := c.Param("post_id")

	var post model.Post
	// 关联查询用户信息
	if err := config.DB.Preload("User").Where("id = ?", postID).First(&post).Error; err != nil {
		logrus.Warnf("查询文章失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	logrus.Infof("查询文章详情成功: %d", post.ID)
	c.JSON(http.StatusOK, gin.H{"post": post})
}

// UpdatePostRequest 更新文章请求参数
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=200"`
	Content string `json:"content" binding:"omitempty,min=10"`
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	postID := c.Param("post_id")
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("更新文章参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 查询文章
	var post model.Post
	if err := config.DB.Preload("User").Where("id = ?", postID).First(&post).Error; err != nil {
		logrus.Warnf("更新文章失败: 文章不存在 %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 验证权限（只有作者能更新）
	userID, _ := c.Get("userID")
	if post.UserID != userID.(uint) {
		logrus.Warnf("用户%d无权限更新文章%d", userID.(uint), post.ID)
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新该文章"})
		return
	}

	// 更新文章信息
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if err := config.DB.Save(&post).Error; err != nil {
		logrus.Errorf("更新文章失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	logrus.Infof("用户%d更新文章成功: %d", userID.(uint), post.ID)
	c.JSON(http.StatusOK, gin.H{"message": "文章更新成功", "post": post})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	postID := c.Param("post_id")

	// 查询文章
	var post model.Post
	if err := config.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		logrus.Warnf("删除文章失败: 文章不存在 %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 验证权限（只有作者能删除）
	userID, _ := c.Get("userID")
	if post.UserID != userID.(uint) {
		logrus.Warnf("用户%d无权限删除文章%d", userID.(uint), post.ID)
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除该文章"})
		return
	}

	// 删除文章（物理删除，也可改为软删除） TODO 未删除评论相关数据
	if err := config.DB.Delete(&post).Error; err != nil {
		logrus.Errorf("删除文章失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	logrus.Infof("用户%d删除文章成功: %d", userID.(uint), post.ID)
	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
