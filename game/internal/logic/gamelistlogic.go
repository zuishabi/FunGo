// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fungo/game/internal/svc"
	"fungo/game/internal/types"
	"fungo/game/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameListLogic {
	return &GameListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameListLogic) GameList(req *types.GameListReq) (resp *types.GemeListRsp, err error) {
	infos := make([]model.GameInfo, 0)
	l.svcCtx.Db.Select("id", "title", "updated_at", "author", "cover", "play_time", "download_file").
		Offset((req.Page - 1) * 10).Limit(10).Find(&infos)
	games := make([]types.GameInfo, len(infos))
	for i, v := range infos {
		games[i].ID = v.ID
		games[i].Author = v.Author
		games[i].Title = v.Title
		games[i].Cover = v.Cover
		games[i].UpdatedAt = v.UpdatedAt.Format("2006.01.02 15:04")
		games[i].CanDownload = v.DownloadFile != ""
		games[i].CanOnlinePlaying = v.CanOnlinePlaying
		games[i].PlayTime = v.PlayTime
	}
	return &types.GemeListRsp{Games: games}, nil
}
