// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"fungo/game/internal/config"
	"fungo/game/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	Db       *gorm.DB
	RedisCli *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		fmt.Println("打开数据库失败")
		panic(err)
	}
	_ = db.AutoMigrate(&model.GameInfo{})

	// 打开redis
	cli := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
	})
	pong, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("连接redis成功：", pong)

	svcCtx := &ServiceContext{
		Config:   c,
		Db:       db,
		RedisCli: cli,
	}

	fmt.Println("启动聚合协程")
	go AggregateOnlinePlayTime(svcCtx)

	return svcCtx
}

func AggregateOnlinePlayTime(svcCTX *ServiceContext) {
	ticker := time.NewTicker(5)
	for {
		select {
		case <-ticker.C:
			res, err := svcCTX.RedisCli.HGetAll(context.Background(), "online-game-play-time").Result()
			if err != nil {
				continue
			}
			for i, v := range res {
				id, err := strconv.Atoi(i)
				if err != nil {
					continue
				}
				// 将数据保存到数据库中
				if err := svcCTX.Db.Model(&model.GameInfo{}).Where("id = ?", id).
					UpdateColumn("play_time", gorm.Expr("play_time + ?", v)).Error; err != nil {
					continue
				}
			}
		}
	}
}
