// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/articles/internal/svc"
	"fungo/articles/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserArticleListLogic {
	return &UserArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserArticleListLogic) UserArticleList(req *types.UserArticleListReq) (resp *types.UserArticleListRsp, err error) {
	// todo: add your logic here and delete this line

	return
}
