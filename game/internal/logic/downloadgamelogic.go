// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadGameLogic {
	return &DownloadGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadGameLogic) DownloadGame(req *types.GameFileReq) error {
	// todo: add your logic here and delete this line

	return nil
}
