// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/animate/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAnimateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAnimateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAnimateLogic {
	return &UploadAnimateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAnimateLogic) UploadAnimate() error {
	// todo: add your logic here and delete this line

	return nil
}
