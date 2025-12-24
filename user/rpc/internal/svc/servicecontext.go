package svc

import (
	"context"
	"fmt"
	"fungo/user/model"
	"fungo/user/rpc/internal/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	RedisCli *redis.Client
	Db       *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("连接到数据库成功")
	}
	_ = db.AutoMigrate(model.User{})

	// 初始化 redis
	cli := redis.NewClient(&redis.Options{
		Addr:     c.RedisCli.Addr,
		Password: c.RedisCli.Password,
	})
	pong, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("连接redis成功：", pong)

	return &ServiceContext{
		Config:   c,
		RedisCli: cli,
		Db:       db,
	}
}
