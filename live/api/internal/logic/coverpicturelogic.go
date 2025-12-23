// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/live/api/internal/svc"
	"fungo/live/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoverPictureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoverPictureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoverPictureLogic {
	return &CoverPictureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoverPictureLogic) CoverPicture(req *types.CoverPictureReq) error {
	return nil
}
