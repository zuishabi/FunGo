// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"fungo/live/api/internal/config"
	"fungo/live/api/internal/types"
	"fungo/live/model"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config           config.Config
	Db               *gorm.DB
	RedisClient      *redis.Client
	BulletChatServer *BulletChatServerType
}

type BulletChatServerType struct {
	lock  sync.Mutex
	rooms map[uint64]map[chan *types.BulletChatMessageRsp]struct{}
}

func (b *BulletChatServerType) CreateConn(roomID uint64) chan *types.BulletChatMessageRsp {
	b.lock.Lock()
	bulletChatMessageChan := make(chan *types.BulletChatMessageRsp, 5)
	if _, ok := b.rooms[roomID]; !ok {
		b.rooms[roomID] = make(map[chan *types.BulletChatMessageRsp]struct{})
	}
	b.rooms[roomID][bulletChatMessageChan] = struct{}{}
	b.lock.Unlock()
	return bulletChatMessageChan
}

func (b *BulletChatServerType) DeleteConn(roomID uint64, conn chan *types.BulletChatMessageRsp) {
	b.lock.Lock()
	delete(b.rooms[roomID], conn)
	b.lock.Unlock()
}

func (b *BulletChatServerType) SendBulletChatMessage(roomID uint64, message string, userName string) {
	b.lock.Lock()
	for i, _ := range b.rooms[roomID] {
		i <- &types.BulletChatMessageRsp{
			Content:  message,
			UserName: userName,
		}
	}
	b.lock.Unlock()
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 在这里连接到数据库
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("连接到数据库成功")
	}
	_ = db.AutoMigrate(model.RoomInfo{})
	_ = db.AutoMigrate(model.UserRoom{})

	// 初始化 redis
	cli := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
	})
	pong, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("连接redis成功：", pong)

	return &ServiceContext{
		Config:           c,
		Db:               db,
		RedisClient:      cli,
		BulletChatServer: &BulletChatServerType{rooms: make(map[uint64]map[chan *types.BulletChatMessageRsp]struct{})},
	}
}
