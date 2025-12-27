// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/user/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditUserCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserCoverLogic {
	return &EditUserCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditUserCoverLogic) EditUserCover() error {
	// todo: add your logic here and delete this line

	return nil
}
