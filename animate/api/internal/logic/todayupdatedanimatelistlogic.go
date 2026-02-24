// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/animate/api/internal/svc"
	"fungo/animate/api/internal/types"
	"fungo/animate/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type TodayUpdatedAnimateListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTodayUpdatedAnimateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TodayUpdatedAnimateListLogic {
	return &TodayUpdatedAnimateListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TodayUpdatedAnimateListLogic) TodayUpdatedAnimateList() (resp *types.GetTodayUpdatedAnimateListRsp, err error) {
	l.svcCtx.AnimateServer.TodayUpdateLock.RLock()
	defer l.svcCtx.AnimateServer.TodayUpdateLock.RUnlock()
	// 从数据库中获取今日更新列表
	var list []model.AnimateList
	l.svcCtx.Db.Where("id in ?", l.svcCtx.TodayUpdate).Find(&list)

	animates := make([]types.AnimateItem, len(list))
	for i, v := range list {
		animates[i].Name = v.Name
		animates[i].ID = v.ID
		animates[i].Description = v.Description
		animates[i].State = v.State
		animates[i].Year = v.Year
	}

	resp = &types.GetTodayUpdatedAnimateListRsp{AnimateItems: animates}

	return
}
