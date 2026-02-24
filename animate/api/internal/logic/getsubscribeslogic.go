// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/animate/model"

	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSubscribesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSubscribesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribesLogic {
	return &GetSubscribesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscribesLogic) GetSubscribes() (resp *types.GetSubscribesRsp, err error) {
	uidValue := l.ctx.Value("user_id")
	if uidValue == nil {
		return nil, errors.New("用户未登录")
	}
	uid := uidValue.(uint64)
	target := make([]model.SubscribeFavoriteList, 0)
	l.svcCtx.Db.Select("id").Where("uid = ? AND subscribe = ?", uid, true).Find(&target)
	res := make([]uint64, len(target))
	for i, v := range target {
		res[i] = v.ID
	}
	return &types.GetSubscribesRsp{
		List: res,
	}, nil
}
