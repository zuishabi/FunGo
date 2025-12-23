// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/live/model"

	"fungo/live/api/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLiveRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLiveRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLiveRoomLogic {
	return &CreateLiveRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLiveRoomLogic) CreateLiveRoom() error {
	// 首先检查直播间是否被创建
	uid := l.ctx.Value("user_id").(uint64)
	if l.svcCtx.Db.Where("uid = ?", uid).First(&model.UserRoom{}).Error == nil {
		return errors.New("直播间已经被创建")
	}
	l.svcCtx.Db.Create(&model.UserRoom{
		UID:    uid,
		RoomID: uid,
	})

	return nil
}
