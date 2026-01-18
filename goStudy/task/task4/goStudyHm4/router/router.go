package router

import (
	"goStudyHm4/controller"
	"goStudyHm4/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 公共路由（无需认证）
	public := r.Group("/api")
	{
		// 认证相关
		public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)

		// 文章相关（查询）
		public.GET("/posts", controller.GetPosts)
		// 嵌套路由组：/api/posts/:post_id
		postGroup := public.Group("/posts/:post_id")
		{
			postGroup.GET("", controller.GetPost)              // 单个文章详情（路径：/api/posts/:post_id）
			postGroup.GET("/comments", controller.GetComments) // 文章评论（路径：/api/posts/:post_id/comments）
		}
	}

	// 私有路由（需认证）
	private := r.Group("/api")
	private.Use(middleware.AuthRequired()) // 应用认证中间件
	{
		// 文章相关（创建、更新、删除）
		private.POST("/posts", controller.CreatePost)
		private.PUT("/posts/:post_id", controller.UpdatePost)
		private.DELETE("/posts/:post_id", controller.DeletePost)

		// 评论相关（创建）
		private.POST("/posts/:post_id/comments", controller.CreateComment)
	}

	return r
}
