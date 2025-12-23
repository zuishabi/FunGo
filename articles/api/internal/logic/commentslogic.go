// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"
	"fungo/articles/model"

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
	comments := make([]model.Comment, 0)
	l.svcCtx.Db.Where("article_id = ?", req.ID).Find(&comments)
	res := make([]types.CommentInfo, len(comments))

	for i, v := range comments {
		res[i].UID = v.UID
		res[i].Content = v.Content
		res[i].CreatedAt = v.CreatedAt.Format("2006.01.02 15:04")
		res[i].ID = v.ID
		res[i].Parent = v.Parent
		res[i].PParent = v.PParent
	}
	return &types.CommentsRsp{Comments: res}, nil
}
