// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/live/model"
	"fungo/user/rpc/user"
	"strconv"

	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLiveRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLiveRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLiveRoomLogic {
	return &UserLiveRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLiveRoomLogic) UserLiveRoom(req *types.UserLiveRoomReq) (resp *types.UserLiveRoomRsp, err error) {
	// 检查用户直播间的信息
	if req.UID == 0 {
		// 尝试通过jwt获得用户id
		ok := true
		req.UID, ok = l.ctx.Value("user_id").(uint64)
		if !ok {
			return nil, errors.New("解析用户id错误")
		}
	}

	room := &model.UserRoom{}
	if err := l.svcCtx.Db.Where("uid = ?", req.UID).First(room).Error; err != nil {
		return &types.UserLiveRoomRsp{
			State: 0,
		}, nil
	}
	// 从数据库中获取数据
	roomInfo := model.RoomInfo{}
	l.svcCtx.Db.Where("room_id = ?", room.RoomID).First(&roomInfo)
	userName, _ := l.svcCtx.UserRPC.GetUserInfo(context.Background(), &user.UserInfoReq{Uid: req.UID})
	LiveInfo := types.LiveInfo{
		UserID:           req.UID,
		UserName:         userName.UserName,
		RoomID:           roomInfo.RoomID,
		Title:            roomInfo.Title,
		CurrentPeopleNum: 0, // 记录当前的在线人数
		Description:      roomInfo.Description,
		Cover:            roomInfo.Cover,
	}
	// 检查redis中是否存在当前直播间
	_, err = l.svcCtx.RedisClient.ZScore(context.Background(), "live-room-list", strconv.Itoa(int(room.RoomID))).Result()
	if err != nil {
		return &types.UserLiveRoomRsp{
			State:    1,
			LiveInfo: LiveInfo,
		}, nil
	}
	return &types.UserLiveRoomRsp{
		State:    2,
		LiveInfo: LiveInfo,
	}, nil
}
