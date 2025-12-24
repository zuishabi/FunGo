// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"fmt"
	"fungo/user/api/internal/config"
	"fungo/user/model"

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
	} else {
		fmt.Println("连接到数据库成功")
	}
	_ = db.AutoMigrate(model.User{})
	return &ServiceContext{
		Config: c,
		Db:     db,
	}
}
