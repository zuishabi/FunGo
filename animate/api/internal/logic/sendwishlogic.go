// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"
	"fungo/animate/model"
	"fungo/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendWishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendWishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendWishLogic {
	return &SendWishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendWishLogic) SendWish(req *types.SendWishReq) (*types.SendWishRsp, error) {
	uid, err := common.CheckLogin(l.ctx)
	if err != nil {
		return nil, err
	}
	// 将许愿数据保存到数据库种
	wish := &model.WishList{
		Content: req.Content,
		UID:     uid,
	}
	err = l.svcCtx.Db.Create(wish).Error
	return &types.SendWishRsp{Wish: types.WishInfo{
		ID:        wish.ID,
		Content:   wish.Content,
		CreatedAt: wish.CreatedAt.Format("2006.01.02"),
	}}, err
}
