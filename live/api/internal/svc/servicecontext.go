// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"fungo/live/api/internal/config"
	"fungo/live/api/internal/types"
	"fungo/live/model"
	"fungo/user/rpc/userclient"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config           config.Config
	Db               *gorm.DB
	RedisClient      *redis.Client
	BulletChatServer *BulletChatServerType
	UserRPC          userclient.User
}

type BulletChatServerType struct {
	lock  sync.RWMutex
	rooms map[uint64]*Room
}

type Room struct {
	lock  sync.RWMutex
	users map[chan *types.BulletChatMessageRsp]struct{}
}

func (b *BulletChatServerType) CreateConn(roomID uint64) chan *types.BulletChatMessageRsp {
	b.lock.RLock()
	defer b.lock.RUnlock()
	bulletChatMessageChan := make(chan *types.BulletChatMessageRsp, 5)
	room, ok := b.rooms[roomID]
	if !ok {
		// 在这里创建房间
		b.rooms[roomID] = &Room{
			users: make(map[chan *types.BulletChatMessageRsp]struct{}),
		}
		room = b.rooms[roomID]
	}
	room.lock.Lock()
	defer room.lock.Unlock()
	room.users[bulletChatMessageChan] = struct{}{}
	return bulletChatMessageChan
}

func (b *BulletChatServerType) DeleteConn(roomID uint64, conn chan *types.BulletChatMessageRsp) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	room := b.rooms[roomID]
	room.lock.Lock()
	defer room.lock.Unlock()
	delete(room.users, conn)
}

func (b *BulletChatServerType) SendBulletChatMessage(roomID uint64, message string, userName string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	room, ok := b.rooms[roomID]
	if !ok {
		// 在这里创建房间
		b.rooms[roomID] = &Room{
			users: make(map[chan *types.BulletChatMessageRsp]struct{}),
		}
		room = b.rooms[roomID]
	}
	room.lock.RLock()
	defer room.lock.RUnlock()
	for i, _ := range room.users {
		i <- &types.BulletChatMessageRsp{
			Content:  message,
			UserName: userName,
		}
	}
}

func (b *BulletChatServerType) DeleteRoom(roomID uint64) {
	b.lock.Lock()
	defer b.lock.Unlock()
	delete(b.rooms, roomID)
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
		BulletChatServer: &BulletChatServerType{rooms: make(map[uint64]*Room)},
		UserRPC:          userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
