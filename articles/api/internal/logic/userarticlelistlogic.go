// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"
	"fungo/articles/model"

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
	// 获取特定用户的所有文章，一页10篇
	if req.UID == 0 {
		// 从jwt中尝试获取用户id
		ok := true
		req.UID, ok = l.ctx.Value("user_id").(uint64)
		if !ok {
			return nil, errors.New("解析用户id失败")
		}
	}
	articles := make([]model.Article, 0)
	l.svcCtx.Db.Where("author = ?", req.UID).Offset((req.Page - 1) * 10).Limit(10).Find(&articles)
	articleInfos := make([]types.ArticleInfo, len(articles))
	for i, v := range articles {
		var urls model.Pictures
		_ = json.Unmarshal([]byte(v.Pictures), &urls)
		var pic string
		if len(urls.Name) > 0 {
			pic = urls.Name[0]
		}
		articleInfos[i].Title = v.Title
		articleInfos[i].ID = v.ID
		articleInfos[i].Author = v.Author
		articleInfos[i].CommentNum = v.CommentNum
		articleInfos[i].LikeNum = v.LikeNum
		articleInfos[i].LookNum = v.LookNum
		articleInfos[i].Time = v.CreatedAt.Format("2006.01.02 15:04")
		articleInfos[i].Section = v.Section
		articleInfos[i].Summary = truncateSummary(v.Content, 20)
		articleInfos[i].PictureUrl = pic
	}
	return &types.UserArticleListRsp{Articles: articleInfos}, nil
}
