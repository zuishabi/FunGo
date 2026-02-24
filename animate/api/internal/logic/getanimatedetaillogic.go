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

type GetAnimateDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAnimateDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnimateDetailLogic {
	return &GetAnimateDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAnimateDetailLogic) GetAnimateDetail(req *types.GetAnimateDetailReq) (resp *types.GetAnimateDetailRsp, err error) {
	// 从数据库中读取数据
	item := model.AnimateList{}
	l.svcCtx.Db.Where("id = ?", req.ID).First(&item)
	resp = &types.GetAnimateDetailRsp{
		Item: types.AnimateItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Year:        item.Year,
			Tags:        nil,
			State:       item.State,
			Num:         item.Num,
		},
	}
	return
}
