// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"

	"fungo/game/internal/svc"
	"fungo/game/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlineGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOnlineGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineGameLogic {
	return &OnlineGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnlineGameLogic) OnlineGame(req *types.OnlineGameReq) error {
	// 在这里将游玩时间记录到redis中进行聚合
	return l.svcCtx.RedisCli.HIncrBy(l.ctx, "online-game-play-time", strconv.Itoa(int(req.ID)), 1).Err()
}
