// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"fungo/live/model"
	"strconv"

	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LiveRoomStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLiveRoomStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveRoomStateLogic {
	return &LiveRoomStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveRoomStateLogic) LiveRoomState() (resp *types.LiveRoomStateRsp, err error) {
	uid := l.ctx.Value("user_id").(uint64)
	// 检查直播间是否创建
	roomInfo := model.UserRoom{}
	if err := l.svcCtx.Db.Where("uid = ?", uid).First(&roomInfo).Error; err != nil {
		return &types.LiveRoomStateRsp{
			State: 1,
		}, nil
	}

	// 检查是否已经开播
	_, err = l.svcCtx.RedisClient.ZScore(context.Background(), "live-room-list", fmt.Sprintf("%d", roomInfo.RoomID)).Result()
	if err == nil {
		// err 为 nil，表示成员存在，直播间已开启
		key, _ := l.svcCtx.RedisClient.HGet(context.Background(), "room-key", strconv.Itoa(int(roomInfo.RoomID))).Result()
		return &types.LiveRoomStateRsp{
			State:  3,
			IP:     "rtmp://112.17.30.188:30003/live",
			Key:    key,
			RoomID: roomInfo.RoomID,
		}, nil
	}

	return &types.LiveRoomStateRsp{
		State:  2,
		RoomID: roomInfo.RoomID,
	}, nil
}
