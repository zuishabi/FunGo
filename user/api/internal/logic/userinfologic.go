// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"
	"fungo/user/model"

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
	if req.UID == 0 {
		ok := true
		req.UID, ok = l.ctx.Value("user_id").(uint64)
		if !ok {
			return nil, errors.New("解析用户id错误")
		}
	}
	user := &model.User{}
	l.svcCtx.Db.Where("id = ?", req.UID).First(&user)

	return &types.UserInfoRsp{
		UserName: user.UserName,
		UserID:   uint64(user.ID),
	}, nil
}
