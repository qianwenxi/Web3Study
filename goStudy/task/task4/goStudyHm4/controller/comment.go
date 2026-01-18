package controller

import (
	"net/http"

	"goStudyHm4/config"
	"goStudyHm4/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreateCommentRequest 创建评论请求参数
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	postID := c.Param("post_id")
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("创建评论参数错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 验证文章是否存在
	var post model.Post
	if err := config.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		logrus.Warnf("创建评论失败: 文章不存在 %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")

	// 创建评论
	comment := model.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  post.ID,
	}
	if err := config.DB.Create(&comment).Error; err != nil {
		logrus.Errorf("创建评论失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	// 关联用户信息返回
	config.DB.Preload("User").Preload("Post").First(&comment)
	logrus.Infof("用户%d为文章%d创建评论成功", userID.(uint), post.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "评论创建成功", "comment": comment})
}

// GetComments 获取文章的所有评论
func GetComments(c *gin.Context) {
	postID := c.Param("post_id")

	// 验证文章是否存在
	var post model.Post
	if err := config.DB.Where("id = ?", postID).First(&post).Error; err != nil {
		logrus.Warnf("查询评论失败: 文章不存在 %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 查询该文章的所有评论（关联用户信息）
	var comments []model.Comment
	if err := config.DB.Preload("User").Preload("Post").Where("post_id = ?", postID).Order("created_at desc").Find(&comments).Error; err != nil {
		logrus.Errorf("查询评论失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	logrus.Infof("查询文章%d的评论成功，共%d条", post.ID, len(comments))
	c.JSON(http.StatusOK, gin.H{"comments": comments, "count": len(comments)})
}
