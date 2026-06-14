// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fungo/animate/model"

	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAnimateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchAnimateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAnimateLogic {
	return &SearchAnimateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchAnimateLogic) SearchAnimate(req *types.SearchAnimateReq) (resp *types.SearchAnimateRsp, err error) {
	// 查询数据库相关关键字
	items := make([]model.AnimateList, 0)
	l.svcCtx.Db.Where("name LIKE ?", "%"+req.Key+"%").Find(&items)

	res := make([]types.AnimateItem, len(items))
	for i, v := range items {
		res[i] = types.AnimateItem{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Year:        v.Year,
			Tags:        nil,
			State:       v.State,
			Num:         v.Num,
		}
	}

	return &types.SearchAnimateRsp{Items: res}, nil
}
