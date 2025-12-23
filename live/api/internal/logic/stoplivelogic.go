// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"fungo/live/api/internal/svc"
	"fungo/live/model"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type StopLiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStopLiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopLiveLogic {
	return &StopLiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopLiveLogic) StopLive() error {
	// 查询用户的房间号
	uid := l.ctx.Value("user_id").(uint64)
	room := model.UserRoom{}
	if err := l.svcCtx.Db.Where("uid = ?", uid).First(&room).Error; err != nil {
		return errors.New("未创建直播间")
	}

	// 清除redis中的对应数据
	l.svcCtx.RedisClient.Del(context.Background(), fmt.Sprintf("live-room-%d", room.RoomID))
	key := l.svcCtx.RedisClient.HGet(context.Background(), "room-key", strconv.Itoa(int(room.RoomID))).String()
	l.svcCtx.RedisClient.HDel(context.Background(), "room-key", strconv.Itoa(int(room.RoomID)))
	l.svcCtx.RedisClient.SRem(context.Background(), "key-set", key)
	l.svcCtx.RedisClient.ZRem(context.Background(), "live-room-list", room.RoomID)

	return nil
}
