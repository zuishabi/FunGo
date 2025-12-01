// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/articles/internal/svc"
	"fungo/articles/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCommentsLogic {
	return &SendCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCommentsLogic) SendComments(req *types.SendCommentsReq) error {
	// todo: add your logic here and delete this line

	return nil
}
