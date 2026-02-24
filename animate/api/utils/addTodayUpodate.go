package main

import (
	"fmt"
	"fungo/animate/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:861214959@tcp(127.0.0.1:3306)/FunGo?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("连接到数据库成功")
	db.Create(&model.TodayUpdateList{
		ID: 1,
	})
	db.Create(&model.TodayUpdateList{
		ID: 2,
	})
}
