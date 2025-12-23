// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/live/api/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditLiveRoomReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditLiveRoomReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditLiveRoomReqLogic {
	return &EditLiveRoomReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditLiveRoomReqLogic) EditLiveRoomReq() error {
	return nil
}
