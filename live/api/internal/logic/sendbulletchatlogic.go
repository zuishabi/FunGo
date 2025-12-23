// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendBulletChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendBulletChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendBulletChatLogic {
	return &SendBulletChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SendBulletChat 发送弹幕
func (l *SendBulletChatLogic) SendBulletChat(req *types.SendBulletChat) error {
	l.svcCtx.BulletChatServer.SendBulletChatMessage(req.RoomID, req.Content, req.UserName)
	return nil
}
