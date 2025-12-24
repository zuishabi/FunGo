// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fungo/user/api/internal/svc"
	"fungo/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelfInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelfInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelfInfoLogic {
	return &SelfInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelfInfoLogic) SelfInfo() (resp *types.SelfInfoRsp, err error) {
	uid, _ := l.ctx.Value("user_id").(json.Number).Int64()
	userName := l.ctx.Value("user_name").(string)
	return &types.SelfInfoRsp{
		UserName: userName,
		UserID:   uint64(uid),
	}, nil
}
