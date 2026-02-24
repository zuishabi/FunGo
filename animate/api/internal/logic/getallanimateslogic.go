// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"fungo/animate/model"

	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllAnimatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllAnimatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAnimatesLogic {
	return &GetAllAnimatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllAnimatesLogic) GetAllAnimates(req *types.GetAllAnimtesReq) (resp *types.GetAllAnimatesRsp, err error) {
	targets := make([]model.AnimateList, 0)
	var maxPage int
	var maxCount int64
	if req.Page == -1 {
		l.svcCtx.Db.Select("id").Limit(10).Offset(0).Order("id desc").Find(&targets)
	} else {
		l.svcCtx.Db.Select("id").Limit(15).Offset((req.Page - 1) * 15).Order("id desc").Find(&targets)
		l.svcCtx.Db.Model(&model.AnimateList{}).Count(&maxCount)
	}
	ids := make([]uint64, len(targets))
	for i, v := range targets {
		ids[i] = v.ID
	}
	if maxCount%15 > 0 {
		maxPage = int(maxCount/15 + 1)
		fmt.Println(1, ":", maxPage)
	} else {
		maxPage = int(maxCount / 15)
		fmt.Println(2, ":", maxPage)
	}

	return &types.GetAllAnimatesRsp{List: ids, MaxPage: maxPage}, nil
}
