package main

import (
	hm3 "GOSTUDY/gorm/hw3"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// 替换成你的实际模块路径
)

// type Parent struct {
// 	ID   int `gorm:"primary_key"`
// 	Name string
// }

// type Child struct {
// 	Parent
// 	Age int
// }

// func InitDB(dst ...interface{}) *gorm.DB {
func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/WEB3STUDY?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	//db.AutoMigrate(dst...)

	return db
}

func main() {
	// db, err := gorm.Open(mysql.Open("root:st123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	// if err != nil {
	// 	panic(err)
	// }

	// lesson01.Run(db)
	// lesson02.Run(db)
	// lesson03.Run(db)
	// lesson03_02.Run(db)
	// lesson03_03.Run(db)
	// lesson03_04.Run(db)
	// lesson04.Run(db)

	db := InitDB()
	hm3.Run(db)
}
