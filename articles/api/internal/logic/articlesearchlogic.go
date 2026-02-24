// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleSearchLogic {
	return &ArticleSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleSearchLogic) ArticleSearch(req *types.ArticleSearchReq) (resp *types.ArticleSearchRsp, err error) {
	// 首先从elasticsearch中获取数据

	return
}
