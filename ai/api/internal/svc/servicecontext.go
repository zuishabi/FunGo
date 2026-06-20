package svc

import (
	"fmt"

	"fungo/ai/api/internal/config"
	"fungo/animate/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("连接到数据库成功")
	db.AutoMigrate(&model.AnimateList{})

	return &ServiceContext{
		Config: c,
		Db:     db,
	}
}
