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

type FavoriteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFavoriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FavoriteLogic {
	return &FavoriteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FavoriteLogic) Favorite(req *types.FavoriteReq) (*types.FavoriteRsp, error) {
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
	target.Favorite = !target.Favorite
	err = l.svcCtx.Db.Save(&target).Error
	res = target.Favorite
	return &types.FavoriteRsp{Favorite: res}, err
}
