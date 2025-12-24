// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"fungo/articles/model"

	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteArticleLogic {
	return &DeleteArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteArticleLogic) DeleteArticle(req *types.DeleteArticle) error {
	// 直接删除数据库
	articleInfo := model.Article{}
	l.svcCtx.Db.Where("id = ?", req.ArticleID).First(&articleInfo)
	l.svcCtx.Db.Where("id = ?", req.ArticleID).Delete(&model.Article{})
	l.svcCtx.Db.Where("article_id = ?", req.ArticleID).Delete(&model.Comment{})

	// 再删除缓存中的信息
	l.svcCtx.RedisClient.ZRem(context.Background(), "section-1-list", req.ArticleID)
	l.svcCtx.RedisClient.ZRem(context.Background(), fmt.Sprintf("section-%d-list", req.ArticleID), req.ArticleID)
	l.svcCtx.RedisClient.ZRem(context.Background(), "hot-article-list", req.ArticleID)
	l.svcCtx.RedisClient.Del(context.Background(), fmt.Sprintf("article-%d", req.ArticleID))

	return nil
}
