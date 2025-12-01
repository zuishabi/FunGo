// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/common/jwts"
	"fungo/user/model"

	"fungo/user/internal/svc"
	"fungo/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRsp, err error) {
	// 检查用户和密码是否正确
	user := model.User{}
	if l.svcCtx.Db.Take(&user, "user_name = ?", req.UserName).Error != nil {
		return nil, errors.New("用户名不存在")
	}
	if !ComparePassword(user.Password, req.Password) {
		return nil, errors.New("密码错误")
	}

	auth := l.svcCtx.Config.Auth
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   uint64(user.ID),
		UserName: user.UserName,
		Role:     1,
	}, auth.AccessSecret, auth.AccessExpire)
	if err != nil {
		return nil, err
	}
	return &types.LoginRsp{Token: token}, nil
}
