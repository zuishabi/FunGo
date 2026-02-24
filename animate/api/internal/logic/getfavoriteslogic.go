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

type GetFavoritesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFavoritesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFavoritesLogic {
	return &GetFavoritesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFavoritesLogic) GetFavorites(req *types.GetFavoritesReq) (resp *types.GetFavoritesRsp, err error) {
	uidValue := l.ctx.Value("user_id")
	if uidValue == nil {
		return nil, errors.New("用户未登录")
	}
	uid := uidValue.(uint64)
	target := make([]model.SubscribeFavoriteList, 0)
	l.svcCtx.Db.Select("id").Limit(15).Offset((req.Page-1)*15).Where("uid = ? AND favorite = ?", uid, true).Find(&target)
	res := make([]uint64, len(target))
	for i, v := range target {
		res[i] = v.ID
	}
	return &types.GetFavoritesRsp{
		List: res,
	}, nil
}
