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

type GetSubscribeInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSubscribeInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubscribeInfoLogic {
	return &GetSubscribeInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubscribeInfoLogic) GetSubscribeInfo() (resp *types.GetSubscribeInfoRsp, err error) {
	userValue := l.ctx.Value("user_id")
	if userValue == nil {
		return nil, errors.New("用户未登录")
	}
	uid := userValue.(uint64)
	target := make([]model.SubscribeFavoriteList, 0)
	l.svcCtx.Db.Select("id").Where("uid = ? AND subscribe = ?", uid, true).Find(&target)
	ids := make([]uint64, len(target))
	for i, v := range target {
		ids[i] = v.ID
	}
	res := make([]model.AnimateUpdateInfo, len(ids))
	l.svcCtx.Db.Where("id IN ?", ids).Find(&res)
	resp = &types.GetSubscribeInfoRsp{Infos: make([]types.SubInfo, len(res))}
	for i, v := range res {
		resp.Infos[i].Name = v.Name
		resp.Infos[i].Description = v.Description
		resp.Infos[i].Time = v.UpdatedAt.Format("2006.01.02")
		resp.Infos[i].ID = v.ID
		resp.Infos[i].Version = v.Version
	}

	return
}
