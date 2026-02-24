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

type QuerySubscribeAndFavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuerySubscribeAndFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerySubscribeAndFavoriteLogic {
	return &QuerySubscribeAndFavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuerySubscribeAndFavoriteLogic) QuerySubscribeAndFavorite(req *types.GetSubscribeAndFavoriteReq) (resp *types.GetSubscribeAndFavoriteRsp, err error) {
	uidValue := l.ctx.Value("user_id")
	if uidValue == nil {
		return nil, errors.New("用户没有登录")
	}
	uid := uidValue.(uint64)
	var target model.SubscribeFavoriteList
	l.svcCtx.Db.Where("id = ? AND uid = ?", req.ID, uid).First(&target)
	resp = &types.GetSubscribeAndFavoriteRsp{
		Subscribe: target.Subscribe,
		Favorite:  target.Favorite,
	}
	return
}
