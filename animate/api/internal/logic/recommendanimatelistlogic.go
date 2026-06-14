// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/animate/model"

	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendAnimateListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendAnimateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendAnimateListLogic {
	return &RecommendAnimateListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendAnimateListLogic) RecommendAnimateList() (resp *types.GetRecommendAnimateListRsp, err error) {
	item := model.AnimateList{}
	l.svcCtx.Db.Where("id = ?", 1).First(&item)
	resp = &types.GetRecommendAnimateListRsp{
		AnimateItems: []types.AnimateItem{
			{ID: item.ID,
				Name:        item.Name,
				Description: item.Description,
				Year:        item.Year,
				Tags:        nil,
				State:       item.State,
				Num:         item.Num,
			},
		},
	}
	return
}
