package model

import (
	"goStudyHm4/config"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;size:50" json:"username"`
	Password string `gorm:"not null;size:100" json:"-"` // json:"-" 避免返回密码
	Email    string `gorm:"unique;not null;size:100" json:"email"`
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title   string `gorm:"not null;size:200" json:"title"`
	Content string `gorm:"not null;type:text" json:"content"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联用户
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null;size:500" json:"content"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	PostID  uint   `gorm:"not null" json:"post_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联用户
	Post    Post   `gorm:"foreignKey:PostID" json:"post,omitempty"` // 关联文章
}

// 自动迁移数据库表
func AutoMigrate() {
	err := config.DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("数据库表迁移失败: " + err.Error())
	}
}
