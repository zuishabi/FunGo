// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BulletChatLogic struct {
	logx.Logger
	ctx         context.Context
	messageChan chan *types.BulletChatMessageRsp
	svcCtx      *svc.ServiceContext
}

func NewBulletChatLogic(ctx context.Context, svcCtx *svc.ServiceContext, messageChan chan *types.BulletChatMessageRsp) *BulletChatLogic {
	return &BulletChatLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		messageChan: messageChan,
		svcCtx:      svcCtx,
	}
}

func (l *BulletChatLogic) BulletChat(client chan<- *types.BulletChatMessageRsp) error {
	// 将消息推送给客户端
	for {
		select {
		case message := <-l.messageChan:
			client <- message
		case <-l.ctx.Done():
			return nil
		}
	}
}
