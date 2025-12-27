// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type LiveRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLiveRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveRoomLogic {
	return &LiveRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveRoomLogic) LiveRoom(req *types.LiveRoomReq) (*types.LiveRoomRsp, error) {
	// 从redis中获取房间信息和key
	roomID := req.RoomID
	result, err := l.svcCtx.RedisClient.HMGet(context.Background(), fmt.Sprintf("live-room-%d", roomID),
		"room_id", "title", "description", "user_id", "cover", "current_people").Result()
	if err != nil || result[1] == nil {
		return nil, errors.New("当前房间不存在")
	}
	rsp := &types.LiveRoomRsp{}
	rsp.RoomID = roomID
	rsp.Title = result[1].(string)
	rsp.Description = result[2].(string)
	uid, _ := strconv.Atoi(result[3].(string))
	rsp.UserID = uint64(uid)
	rsp.Cover = result[4].(string)
	rsp.CurrentPeopleNum, _ = strconv.Atoi(result[5].(string))

	rsp.Key, _ = l.svcCtx.RedisClient.HGet(context.Background(), "room-key", strconv.Itoa(int(roomID))).Result()
	return rsp, nil
}
