// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/user/model"

	"fungo/user/internal/svc"
	"fungo/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoRsp, err error) {
	user := &model.User{}
	l.svcCtx.Db.Where("id = ?", req.UID).First(&user)

	return &types.UserInfoRsp{
		UserName: user.UserName,
		UserID:   uint64(user.ID),
	}, nil
}
