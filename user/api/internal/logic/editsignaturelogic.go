// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/user/model"

	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditSignatureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditSignatureLogic {
	return &EditSignatureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditSignatureLogic) EditSignature(req *types.EditSignatureReq) error {
	userValue := l.ctx.Value("user_id")
	if userValue == nil {
		return errors.New("用户未登录")
	}
	uid := userValue.(uint64)

	l.svcCtx.Db.Model(&model.User{}).Where("id = ?", uid).Update("signature", req.Content)

	return nil
}
