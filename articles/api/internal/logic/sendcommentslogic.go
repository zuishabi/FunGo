// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"
	"fungo/articles/model"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCommentsLogic {
	return &SendCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCommentsLogic) SendComments(req *types.SendCommentsReq) error {
	user, _ := l.ctx.Value("user_id").(uint64)
	comment := model.Comment{
		UID:       user,
		ArticleID: req.ArticleID,
		Parent:    req.Parent,
		PParent:   req.PParent,
		Content:   req.Content,
	}

	// 更新热度排行榜
	l.svcCtx.RedisClient.ZIncrBy(context.Background(), "hot-article-list", 5, strconv.Itoa(int(req.ArticleID)))

	return l.svcCtx.Db.Create(&comment).Error
}
