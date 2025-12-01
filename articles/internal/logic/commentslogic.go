// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/articles/internal/svc"
	"fungo/articles/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentsLogic {
	return &CommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentsLogic) Comments(req *types.CommentsReq) (resp *types.CommentsRsp, err error) {
	// todo: add your logic here and delete this line

	return
}
