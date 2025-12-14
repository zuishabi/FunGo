// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPicturesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPicturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPicturesLogic {
	return &GetPicturesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPicturesLogic) GetPictures(req *types.PictureReq) error {
	// todo: add your logic here and delete this line

	return nil
}
