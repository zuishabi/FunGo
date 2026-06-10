// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fungo/animate/model"
	"math/rand"
	"time"

	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWishLogic {
	return &GetWishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWishLogic) GetWish() (resp *types.GetWishRsp, err error) {
	// 生成一个随机数
	wish := &model.WishList{}
	l.svcCtx.Db.Last(wish)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tryTime := 0
	for tryTime < 5 {
		randomNum := r.Int63n(int64(wish.ID)) + 1
		wish = &model.WishList{}
		if err := l.svcCtx.Db.Where("id = ?", randomNum).First(wish).Error; err != nil {
			tryTime += 1
			continue
		} else {
			break
		}
	}
	if tryTime == 5 {
		return nil, errors.New("请求失败")
	}

	// 构建响应
	resp = &types.GetWishRsp{
		Wish: types.WishInfo{
			ID:        wish.ID,
			Content:   wish.Content,
			CreatedAt: wish.CreatedAt.Format("2006-01-02"),
		},
	}
	return
}
