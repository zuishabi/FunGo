// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"fungo/user/rpc/user"
	"strconv"

	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LiveRoomListReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLiveRoomListReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveRoomListReqLogic {
	return &LiveRoomListReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveRoomListReqLogic) LiveRoomListReq(req *types.LiveRoomListReq) (resp *types.LiveRoomListRsp, err error) {
	res, err := l.svcCtx.RedisClient.ZRevRange(context.Background(), "live-room-list", int64((req.Page-1)*12), int64(req.Page*12)).Result()
	if err != nil {
		return nil, err
	}
	liveInfos := make([]types.LiveInfo, len(res))
	for i, v := range res {
		roomID, _ := strconv.Atoi(v)
		result, _ := l.svcCtx.RedisClient.HMGet(context.Background(), fmt.Sprintf("live-room-%d", roomID),
			"room_id", "title", "description", "user_id", "cover", "current_people").Result()
		liveInfos[i].RoomID = uint64(roomID)
		liveInfos[i].Title = result[1].(string)
		liveInfos[i].Description = result[2].(string)
		uid, _ := strconv.Atoi(result[3].(string))
		liveInfos[i].UserID = uint64(uid)
		liveInfos[i].Cover = result[4].(string)
		userInfo, _ := l.svcCtx.UserRPC.GetUserInfo(context.Background(), &user.UserInfoReq{Uid: liveInfos[i].UserID})
		liveInfos[i].UserName = userInfo.UserName
		liveInfos[i].CurrentPeopleNum, _ = strconv.Atoi(result[5].(string))
	}

	return &types.LiveRoomListRsp{LiveInfos: liveInfos}, nil
}
