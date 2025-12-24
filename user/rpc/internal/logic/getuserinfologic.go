package logic

import (
	"context"
	"errors"
	"fmt"
	"fungo/user/model"
	"fungo/user/rpc/internal/svc"
	"fungo/user/rpc/user"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.UserInfoReq) (*user.UserInfoRsp, error) {
	var userInfo model.User
	rsp := user.UserInfoRsp{Uid: in.Uid}

	key := fmt.Sprintf("user-info-%d", in.Uid)

	vals, err := l.svcCtx.RedisCli.HMGet(context.Background(), key, "user_name", "created_at").Result()
	cacheHit := err == nil && len(vals) == 2 && vals[0] != nil && vals[1] != nil

	if cacheHit {
		rsp.UserName = fmt.Sprint(vals[0])
		rsp.CreatedAt = fmt.Sprint(vals[1])
		return &rsp, nil
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		l.Logger.Errorf("redis HMGet failed, key=%s, err=%v", key, err)
	}

	if dbErr := l.svcCtx.Db.Where("id = ?", in.Uid).First(&userInfo).Error; dbErr != nil {
		return nil, dbErr
	}

	rsp.UserName = userInfo.UserName
	rsp.CreatedAt = userInfo.CreatedAt.Format("2006.01.02 15:04")

	pipe := l.svcCtx.RedisCli.Pipeline()
	pipe.HSet(context.Background(), key, "user_name", rsp.UserName, "created_at", rsp.CreatedAt)
	pipe.Expire(context.Background(), key, 30*time.Minute) // 注释与实际保持一致
	if _, werr := pipe.Exec(context.Background()); werr != nil {
		l.Logger.Errorf("redis write cache failed, key=%s, err=%v", key, werr)
	}
	return &rsp, nil
}
