// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadGameLogic {
	return &UploadGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadGameLogic) UploadGame(req *types.UploadGameReq) error {
	// todo: add your logic here and delete this line

	return nil
}
