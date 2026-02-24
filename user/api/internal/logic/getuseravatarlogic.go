// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAvatarLogic {
	return &GetUserAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserAvatarLogic) GetUserAvatar(req *types.GetUserAvatarReq) error {
	// todo: add your logic here and delete this line

	return nil
}
