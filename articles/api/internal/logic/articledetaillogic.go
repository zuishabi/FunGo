// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"fungo/articles/api/internal/svc"
	"fungo/articles/api/internal/types"
	"fungo/articles/model"
	"fungo/common/bitmap"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ArticleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDetailLogic {
	return &ArticleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDetailLogic) ArticleDetail(req *types.ArticleDetailReq) (resp *types.ArticleDetailRsp, err error) {
	curUserID := l.ctx.Value("user_id")
	var uid uint64
	if curUserID == nil {
		fmt.Println("游客查看")
	} else {
		uid, _ = curUserID.(uint64)
	}

	key := fmt.Sprintf("article-%d", req.ID)
	res, err := l.svcCtx.RedisClient.HMGet(context.Background(), key,
		"id", "author", "time", "title", "section", "comment_num", "look_num", "like_num", "pictures", "content").Result()
	if err != nil {
		return nil, err
	}

	allNil := true
	for _, v := range res {
		if v != nil {
			allNil = false
			break
		}
	}

	safeString := func(idx int) string {
		if idx < 0 || idx >= len(res) {
			return ""
		}
		if res[idx] == nil {
			return ""
		}
		if s, ok := res[idx].(string); ok {
			return s
		}
		if b, ok := res[idx].([]byte); ok {
			return string(b)
		}
		return fmt.Sprint(res[idx])
	}

	safeInt := func(idx int) int {
		s := safeString(idx)
		if s == "" {
			return 0
		}
		i, _ := strconv.Atoi(s)
		return i
	}

	// 获取是否点赞的信息
	var isLike bool = false
	b, err := l.svcCtx.RedisClient.GetBit(context.Background(), fmt.Sprintf("likebitmap-%d", req.ID), int64(uid)).Result()
	if err == nil {
		isLike = b == 1
	} else {
		// 从数据库中获得是否点赞
		article := &model.Article{}
		l.svcCtx.Db.Select("likes_bitmap").Where("id = ?", req.ID).First(article)
		isLike = bitmap.IsBitSet(article.LikesBitmap, uid)
		fmt.Println(isLike)
	}

	if !allNil {
		// 当前键存在
		var pic model.Pictures
		picStr := safeString(8)
		if picStr != "" {
			_ = json.Unmarshal([]byte(picStr), &pic)
		}

		articleID := uint64(safeInt(0))
		authorID := uint64(safeInt(1))
		section := uint32(safeInt(4))
		commentNum := safeInt(5)
		lookNum := safeInt(6)
		likeNum := safeInt(7)

		rsp := &types.ArticleDetailRsp{
			ArticleInfo: types.ArticleInfo{
				ID:         articleID,
				Author:     authorID,
				Time:       safeString(2),
				Title:      safeString(3),
				Section:    section,
				CommentNum: commentNum,
				LookNum:    lookNum,
				LikeNum:    likeNum,
				IfLike:     isLike,
			},
			Pictures: pic.Name,
			Content:  safeString(9),
		}

		l.svcCtx.RedisClient.HIncrBy(context.Background(), key, "look_num", 1)
		l.svcCtx.RedisClient.HIncrBy(context.Background(), "article-look-nums", strconv.Itoa(int(req.ID)), 1)
		return rsp, nil
	}

	// 从mysql中读入
	var article model.Article
	if err := l.svcCtx.Db.Find(&article, "id = ?", req.ID).Error; err != nil {
		return nil, err
	}
	// 让观看数+1
	if err := l.svcCtx.Db.Model(&model.Article{}).
		Where("id = ?", req.ID).
		UpdateColumn("look_num", gorm.Expr("look_num + ?", 1)).Error; err != nil {
		fmt.Println("db increment look_num error:", err)
	} else {
		// 本地同步加一，后续返回或缓存使用这个值
		article.LookNum = article.LookNum + 1
	}

	pictures := model.Pictures{}
	_ = json.Unmarshal([]byte(article.Pictures), &pictures)
	rsp := &types.ArticleDetailRsp{
		ArticleInfo: types.ArticleInfo{
			ID:         article.ID,
			Author:     article.Author,
			Time:       article.CreatedAt.Format("2006.01.02 15:04"),
			Title:      article.Title,
			Section:    article.Section,
			CommentNum: article.CommentNum,
			LookNum:    article.LookNum,
			LikeNum:    article.LikeNum,
			IfLike:     isLike,
		},
		Pictures: pictures.Name,
		Content:  article.Content,
	}
	// 当查看数大于100时才放入redis缓存
	if rsp.LookNum > 100 {
		// 将帖子元信息存储到redis中
		err = l.svcCtx.RedisClient.HMSet(context.Background(), fmt.Sprintf("article-%d", req.ID),
			"id", rsp.ID, "author", rsp.Author, "time", rsp.Time, "title", rsp.Title, "section", rsp.Section, "comment_num", rsp.CommentNum,
			"look_num", rsp.LookNum, "like_num", rsp.LikeNum, "pictures", article.Pictures, "content", rsp.Content).Err()
		if err != nil {
			fmt.Println(err)
		}
	}

	// 设置热度排行
	l.svcCtx.RedisClient.ZIncrBy(context.Background(), "hot-article-list", 1, strconv.Itoa(int(article.ID)))

	return rsp, nil
}
