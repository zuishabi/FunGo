// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"
	"fungo/animate/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWishListLogic {
	return &GetWishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWishListLogic) GetWishList() (resp *types.GetWishListRsp, err error) {
	// 获取最新的5个许愿
	wishes := make([]model.WishList, 0)
	l.svcCtx.Db.Limit(5).Offset(0).Order("id desc").Find(&wishes)
	rsp := make([]types.WishInfo, len(wishes))
	for i, v := range wishes {
		rsp[i].ID = v.ID
		rsp[i].Content = v.Content
		rsp[i].CreatedAt = v.CreatedAt.Format("2006.01.02")
	}

	return &types.GetWishListRsp{Wishes: rsp}, nil
}
