// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/user/model"

	"fungo/user/internal/svc"
	"fungo/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (l *RegisterLogic) Register(req *types.RegisterReq) error {
	// 首先检查当前用户是否已经注册
	if l.svcCtx.Db.Take(&model.User{}, "user_name = ?", req.UserName).Error == nil {
		// 已经注册过了
		return errors.New("当前用户已经注册")
	}
	pwd, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
	l.svcCtx.Db.Create(&model.User{
		UserName: req.UserName,
		Password: pwd,
	})
	return nil
}
