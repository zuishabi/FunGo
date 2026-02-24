// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"
	"fungo/animate/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeLogic {
	return &SubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeLogic) Subscribe(req *types.SubscribeReq) (*types.SubscribeRsp, error) {
	uidValue := l.ctx.Value("user_id")
	if uidValue == nil {
		return nil, errors.New("用户没有登录")
	}
	uid := uidValue.(uint64)
	var target model.SubscribeFavoriteList
	res := false
	err := l.svcCtx.Db.Where("id = ? AND uid = ?", req.Id, uid).First(&target).Error
	if err != nil {
		target.ID = req.Id
		target.UID = uid
	}
	target.Subscribe = !target.Subscribe
	err = l.svcCtx.Db.Save(&target).Error
	res = target.Subscribe
	return &types.SubscribeRsp{Subscribe: res}, err
}
