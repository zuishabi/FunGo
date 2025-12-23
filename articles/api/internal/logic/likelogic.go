// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"
	"fungo/articles/model"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var setlikenum = `
	local id = ARGV[1]
	if not id or id == '' then
	  return redis.error_reply("invalid id")
	end
	
	local key = "article-" .. id
	if redis.call("EXISTS", key) == 1 then
	  -- 对 hash 字段 like_num 原子 +1，返回新值
	  local new = redis.call("HINCRBY", key, "like_num", 1)
	  -- 记录到待落库集合
	  redis.call("HSET", "like-save-set", id,new)
	  return new
	else
	  return redis.error_reply("key not exists: " .. key)
	end
`

var checkiflike = `
	local aid = ARGV[1]
	local uid_s = ARGV[2]
	
	if not aid or aid == '' then
		return redis.error_reply("invalid article id")
	end
	if not uid_s or uid_s == '' then
		return redis.error_reply("invalid user id")
	end
	
	local uid = tonumber(uid_s)
	if not uid then
		return redis.error_reply("invalid user id")
	end
	
	local bitmapKey = "likebitmap-" .. aid
	
	-- 检查 bitmap key 是否存在
	if redis.call("EXISTS", bitmapKey) == 0 then
		return redis.error_reply("bitmap not exists")
	end
	
	-- 返回 0 或 1 表示该 uid 在 bitmap 中是否已置位
	local bit = redis.call("GETBIT", bitmapKey, uid)
	return bit
`

type LikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeLogic) Like(req *types.LikeReq) error {
	uid, _ := l.ctx.Value("user_id").(uint64)
	// 首先检查是否已经点赞，首先尝试获取缓存中的bitmap，如果不存在则从数据库中写入
	res, err := l.svcCtx.RedisClient.Eval(context.Background(), checkiflike, nil, req.ID, uid).Result()
	if err == nil {
		fmt.Println("缓存中有点赞bitmap = ", res)
		// 检查是否已经点赞
		ifLike := false
		switch v := res.(type) {
		case int64:
			ifLike = v == 1
		case string:
			n, _ := strconv.ParseInt(v, 10, 64)
			ifLike = n == 1
		case []byte:
			n, _ := strconv.ParseInt(string(v), 10, 64)
			ifLike = n == 1
		}
		if ifLike {
			return errors.New("您已经点赞过了")
		}
	} else {
		// 从数据库中写入
		fmt.Println("缓存中没有点赞bitmap")
		article := model.Article{}
		l.svcCtx.Db.Select("likes_bitmap").Where("id = ?", req.ID).First(&article)
		l.svcCtx.RedisClient.SetNX(context.Background(), fmt.Sprintf("likebitmap-%d", req.ID), article.LikesBitmap, 0)
		ifLike, _ := l.svcCtx.RedisClient.GetBit(context.Background(), fmt.Sprintf("likebitmap-%d", req.ID), int64(uid)).Result()
		if ifLike == 1 {
			return errors.New("您已经点赞过了")
		}
	}

	// 更新点赞数，首先检查是否存在元数据，如果存在就直接进行加1操作，否则修改数据库的数据
	_, err = l.svcCtx.RedisClient.Eval(context.Background(), setlikenum, nil, fmt.Sprintf("%d", req.ID)).Result()
	if err != nil {
		// 脚本返回的错误或连接错误都会在 err 中体现
		l.svcCtx.Db.Model(&model.Article{}).Where("id = ?", req.ID).UpdateColumn("like_num", gorm.Expr("like_num + ?", 1))
	}

	// 更新点赞bitmap和热度排行榜
	l.svcCtx.RedisClient.SetBit(context.Background(), fmt.Sprintf("likebitmap-%d", req.ID), int64(uid), 1)
	l.svcCtx.RedisClient.SAdd(context.Background(), "likebitmap-save-list", req.ID)
	l.svcCtx.RedisClient.ZIncrBy(context.Background(), "hot-article-list", 2, strconv.Itoa(int(req.ID)))

	return nil
}
