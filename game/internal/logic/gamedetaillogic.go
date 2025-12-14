// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fungo/game/model"

	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GameDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGameDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameDetailLogic {
	return &GameDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameDetailLogic) GameDetail(req *types.GameDetailReq) (resp *types.GameDetailRsp, err error) {
	// 从数据库中查询数据
	gameInfo := model.GameInfo{}
	if err := l.svcCtx.Db.Where("id = ?", req.ID).First(&gameInfo).Error; err != nil {
		return nil, err
	}
	gameInfoRsp := types.GameInfo{
		ID:               req.ID,
		Title:            gameInfo.Title,
		Cover:            gameInfo.Cover,
		Author:           gameInfo.Author,
		UpdatedAt:        gameInfo.UpdatedAt.Format("2006.01.02 15:04"),
		CanDownload:      gameInfo.DownloadFile != "",
		CanOnlinePlaying: gameInfo.CanOnlinePlaying,
		PlayTime:         gameInfo.PlayTime,
	}

	links := make(map[string]string)
	if err := json.Unmarshal([]byte(gameInfo.Links), &links); err != nil {
		return nil, err
	}

	pictures := model.Pictures{}
	if err := json.Unmarshal([]byte(gameInfo.Pictures), &pictures); err != nil {
		return nil, err
	}

	res := &types.GameDetailRsp{
		Info:        gameInfoRsp,
		Pictures:    pictures.Name,
		Description: gameInfo.Description,
		Links:       links,
	}
	return res, nil
}
