package hm3

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"size:50; not null; unique"`               // 用户名（唯一）
	Email        string `gprm:"varchar(100); size:100; not null;unique"` //邮箱
	PasswordHash string `gorm:"size:255; not null"`                      //密码哈希
	PostCount    int    `gorm:"default:0"`                               // 文章数量统计
	Posts        []Post `gorm:"foreignKey:UserID"`                       // 一对多关联，一个用户多篇文章
}

type Post struct {
	gorm.Model
	Title         string    `gorm:"size:100;not null"`   //文章标题
	Content       string    `gorm:"type:text;not null"`  //文章内容
	UserID        uint      `gorm:"not null"`            //关联用户ID
	CommentCount  int       `gorm:"default:0"`           //评论数量统计
	CommentStatus string    `gorm:"size:20;default:'Y'"` //评论状态（Y：有评论 N：无评论）
	Comments      []Comment `gorm:"foreignKey:PostID"`   //一对多关联：一篇文章有多条评论
}

// 创建文章时更新用户文章数
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	var user User
	if err := tx.First(&user, p.UserID).Error; err != nil {
		return err // 用户不存在终止创建
	}
	return tx.Model(&User{}).Where("id=?", p.UserID).Update("post_count", gorm.Expr("post_count + ?", 1)).Error
}

// 删除评论时，更新文章状态
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 查询当前评论对应的文章
	var post Post
	if err := tx.First(&post, c.PostID).Error; err != nil {
		return err
	}

	// 重新计算文章的评论数
	var commentCount int64
	if err := tx.Model(&Comment{}).Where("post_id = ? AND deleted_at is Null", c.PostID).Count(&commentCount).Error; err != nil {
		return err
	}

	// 更新文章的评论数和状态
	updateData := map[string]interface{}{
		"comment_count": commentCount,
	}
	if commentCount == 0 {
		updateData["comment_status"] = "N"
	} else {
		updateData["comment_status"] = "Y"
	}
	return tx.Model(&Post{}).Where("id=?", c.PostID).Updates(updateData).Error
}

type Comment struct {
	gorm.Model
	Content   string    `gorm:"type:text; not null"` //评论内容
	PostID    uint      `gorm:"not null"`            // 关联文章ID外键
	UserID    uint      `gorm:"not null"`            // 评论用户ID 外键 可关联用户模型
	CreatedAt time.Time `gorm:"autoCreateTime"`      // 创建时间（自动填充）
}

func Run(db *gorm.DB) {
	// 自动迁移创建表（会根据模型创建/更新表结构，不会删除已有数据）
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("表创建失败：" + err.Error())
	}

	println("数据库表创建成功！")

	// 添加测试数据
	users := []User{
		{
			Username:     "xiaozhang",
			Email:        "xiaozhang@123.com",
			PasswordHash: "e10adc3949ba59abbe56e057f20f883e", // 123456的md5哈希，仅作测试
		}, {
			Username:     "xiaoli",
			Email:        "xiaoli@123.com",
			PasswordHash: "e10adc3949ba59abbe56e057f20f883e", // 123456的md5哈希，仅作测试
		},
	}
	if err := db.Create(&users).Error; err != nil {
		log.Printf("创建用户失败：%v", err)
	} else {
		log.Printf("成功创建 %d 个测试用户", len(users))
	}

	posts := []Post{
		{Title: "xiaozhang文章一", Content: "我是xiaozhang的文章一的内容", UserID: 1},
		{Title: "xiaozhang文章二", Content: "我是xiaozhang的文章二的内容", UserID: 1},
		{Title: "xiaozhang文章三", Content: "我是xiaozhang的文章三的内容", UserID: 1},
		{Title: "xiaoli文章一", Content: "我是xiaoli的文章一的内容", UserID: 2},
		{Title: "xiaoli文章二", Content: "我是xiaoli的文章二的内容", UserID: 2},
		{Title: "xiaoli文章三", Content: "我是xiaoli的文章三的内容", UserID: 2},
		{Title: "xiaoli文章四", Content: "我是xiaoli的文章四的内容", UserID: 2},
	}

	if err := db.Create(&posts).Error; err != nil {
		log.Printf("创建文章失败：%v", err)
	} else {
		log.Printf("成功创建 %d 篇测试文章", len(posts))
	}

	comments := []Comment{
		{Content: "评论A", PostID: 1, UserID: 1},
		{Content: "评论B", PostID: 1, UserID: 2},
		{Content: "评论C", PostID: 2, UserID: 1},
		{Content: "评论D", PostID: 3, UserID: 1},
		{Content: "评论E", PostID: 3, UserID: 2},
		{Content: "评论F", PostID: 4, UserID: 1},
		{Content: "评论G", PostID: 5, UserID: 1},
		{Content: "评论H", PostID: 6, UserID: 2},
		{Content: "评论I", PostID: 7, UserID: 1},
	}

	if err := db.Create(&comments).Error; err != nil {
		log.Fatalf("创建评论失败: %v", err)
	} else {
		log.Printf("成功创建 %d 条测试评论", len(comments))
	}
	//db.AutoMigrate(&User{})
	// db.AutoMigrate(&Member{})
	// db.AutoMigrate(&Blog{})
	//db.AutoMigrate(&Blog2{})

	// user := &User{}
	// user.MemberNumber.Valid = true
	// db.Create(user)

	// create传指针
	// mem := Member{}
	// db.Create(&mem)
	// fmt.Println(mem.ID)
	// db.Delete(&Member{}, 1)
}
