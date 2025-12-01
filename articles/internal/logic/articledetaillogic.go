// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"fungo/articles/model"
	"strconv"

	"fungo/articles/internal/svc"
	"fungo/articles/internal/types"

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
	// 首先检查redis中是否有对应的数据
	exist, err := l.svcCtx.RedisClient.Exists(context.Background(), fmt.Sprintf("article-%d", req.ID)).Result()
	if err != nil {
		return nil, err
	}
	if exist > 0 {
		res, err := l.svcCtx.RedisClient.HMGet(context.Background(), fmt.Sprintf("article-%d", req.ID),
			"id", "author", "time", "title", "section", "comment_num", "look_num", "like_num", "pictures", "content").Result()
		if err != nil {
			fmt.Println(err)
		} else {

			var pic model.Pictures
			_ = json.Unmarshal([]byte(res[8].(string)), &pic)
			id, _ := strconv.Atoi(res[0].(string))
			author, _ := strconv.Atoi(res[1].(string))
			section, _ := strconv.Atoi(res[4].(string))
			commentNum, _ := strconv.Atoi(res[5].(string))
			lookNum, _ := strconv.Atoi(res[6].(string))
			likeNUm, _ := strconv.Atoi(res[7].(string))

			rsp := &types.ArticleDetailRsp{
				ArticleInfo: types.ArticleInfo{
					ID:         uint64(id),
					Author:     uint64(author),
					Time:       res[2].(string),
					Title:      res[3].(string),
					Section:    uint32(section),
					CommentNum: commentNum,
					LookNum:    lookNum,
					LikeNum:    likeNUm,
				},
				Pictures: pic.Name,
				Content:  res[9].(string),
			}
			l.svcCtx.RedisClient.HIncrBy(context.Background(), fmt.Sprintf("article-%d", req.ID), "look_num", 1)
			l.svcCtx.RedisClient.HIncrBy(context.Background(), "article-look-nums", strconv.Itoa(int(req.ID)), 1)
			return rsp, err
		}
	}

	var article model.Article
	if err := l.svcCtx.Db.Find(&article, "id = ?", req.ID).Error; err != nil {
		return nil, err
	}
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
		},
		Pictures: pictures.Name,
		Content:  article.Content,
	}
	// 当查看数大于100时才放入redis缓存
	if rsp.LookNum > 100 {
		err = l.svcCtx.RedisClient.HMSet(context.Background(), fmt.Sprintf("article-%d", req.ID),
			"id", rsp.ID, "author", rsp.Author, "time", rsp.Time, "title", rsp.Title, "section", rsp.Section, "comment_num", rsp.CommentNum,
			"look_num", rsp.LookNum, "like_num", rsp.LikeNum, "pictures", article.Pictures, "content", rsp.Content).Err()
		if err != nil {
			fmt.Println(err)
		}
	}
	return rsp, nil
}
