// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/common/jwts"

	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OpLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpLoginLogic {
	return &OpLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OpLoginLogic) OpLogin(req *types.OpLoginReq) (resp *types.LoginRsp, err error) {
	if req.UserName != "admin" || req.Password != "861214959" {
		return nil, errors.New("管理员账号错误")
	}
	auth := l.svcCtx.Config.Auth
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   uint64(0),
		UserName: req.UserName,
		Role:     2,
	}, auth.AccessSecret, auth.AccessExpire)

	return &types.LoginRsp{Token: token}, nil
}
