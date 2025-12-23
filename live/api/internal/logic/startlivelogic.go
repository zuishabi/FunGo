// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"fungo/live/api/internal/svc"
	"fungo/live/model"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

var generateKey = `
	local key = KEYS[1]
	local room_id = ARGV[1]
	
	-- 检查 'key-set' 集合中是否存在该密钥
	if redis.call('SISMEMBER', 'key-set', key) == 1 then
	  -- 如果存在，返回一个包含错误信息的 table，便于客户端处理
	  return {err = "key already exists"}
	end
	
	-- 如果不存在，将密钥添加到 'key-set' 集合
	redis.call('SADD', 'key-set', key)
	-- 将房间号和密钥的映射关系写入 'room-key' 哈希表
	redis.call('HSET', 'room-key', room_id, key)
	
	-- 返回 0 表示成功
	return 0
`

type StartLiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartLiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartLiveLogic {
	return &StartLiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (l *StartLiveLogic) StartLive() error {
	// 首先获得用户的房间号
	uid := l.ctx.Value("user_id").(uint64)
	room := model.UserRoom{}
	if err := l.svcCtx.Db.Where("uid = ?", uid).First(&room).Error; err != nil {
		return errors.New("未创建直播间")
	}
	// 检查当前房间是否已经开启
	_, err := l.svcCtx.RedisClient.ZScore(context.Background(), "live-room-list", fmt.Sprintf("%d", room.RoomID)).Result()
	if err == nil {
		// err 为 nil，表示成员存在，直播间已开启
		return errors.New("直播间已经开启，请勿重复操作")
	}

	// 接着生成一个密钥
	key := randStr(6)
	// 在redis中检查是否有这个密钥，如果有则重新生成一个，否则写入hash中
	for l.svcCtx.RedisClient.Eval(context.Background(), generateKey, []string{key}, room.RoomID).Err() != nil {
		key = randStr(6)
	}

	// 将房间信息再写入redis中
	roomInfo := model.RoomInfo{}
	l.svcCtx.Db.Where("room_id = ?", room.RoomID).First(&roomInfo)
	timeScore := float64(time.Now().UnixNano()) / 1e6
	l.svcCtx.RedisClient.ZAdd(context.Background(), "live-room-list", redis.Z{
		Score:  timeScore,
		Member: room.RoomID,
	})
	l.svcCtx.RedisClient.HMSet(context.Background(), fmt.Sprintf("live-room-%d", room.RoomID),
		"room_id", roomInfo.RoomID,
		"title", roomInfo.Title,
		"description", roomInfo.Description,
		"user_id", uid,
		"cover", roomInfo.Cover,
	)

	return nil
}
