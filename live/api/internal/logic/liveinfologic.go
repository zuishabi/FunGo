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

type LiveInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLiveInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveInfoLogic {
	return &LiveInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveInfoLogic) LiveInfo(req *types.LiveInfoReq) (resp *types.LiveInfo, err error) {
	roomID := req.RoomID
	rsp := &types.LiveInfo{}
	rsp.RoomID = roomID
	// 首先检查是否已经开播，如果已经开播则从redis中获取数据，否则从mysql中获取数据
	exists, _ := l.svcCtx.RedisClient.Exists(context.Background(), fmt.Sprintf("live-room-%d", roomID)).Result()
	if exists == 1 {
		result, _ := l.svcCtx.RedisClient.HMGet(context.Background(), fmt.Sprintf("live-room-%d", roomID),
			"room_id", "title", "description", "user_id", "cover").Result()
		rsp.Title = result[1].(string)
		rsp.Description = result[2].(string)
		uid, _ := strconv.Atoi(result[3].(string))
		rsp.UserID = uint64(uid)
		rsp.Cover = result[4].(string)
	} else {
		// 从数据库中获取数据
		roomInfo := model.RoomInfo{}
		l.svcCtx.Db.Where("room_id = ?", roomID).First(&roomInfo)
		rsp.Title = roomInfo.Title
		rsp.Cover = roomInfo.Cover
		rsp.Description = roomInfo.Description
	}

	return rsp, nil
}
