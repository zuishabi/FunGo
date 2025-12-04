package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"fungo/articles/internal/svc"
	"fungo/articles/internal/types"
	"fungo/articles/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func truncateSummary(s string, limit int) string {
	s = strings.TrimSpace(s)
	if limit <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= limit {
		return string(r)
	}
	return string(r[:limit]) + "..."
}

func NewArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleListLogic {
	return &ArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleListLogic) ArticleList(req *types.ArticleListReq) (resp *types.ArticleListRsp, err error) {
	var ids []int64
	if req.Section == 0 {
		// 显示热门帖子（此处略）
		res, err := l.svcCtx.RedisClient.ZRange(context.Background(), "hot-article-list", 0, 10).Result()
		if err != nil {
			return nil, err
		}
		ids = make([]int64, len(res))
		for i, v := range res {
			//解析出文章id
			articleID, _ := strconv.Atoi(v)
			ids[i] = int64(articleID)
		}
	} else {
		pageSize := 10
		if req.Page <= 0 {
			req.Page = 1
		}
		start := (req.Page - 1) * pageSize
		key := fmt.Sprintf("section-%d-list", req.Section)

		values, err := l.svcCtx.RedisClient.ZRevRange(l.ctx, key, int64(start), int64(start+pageSize-1)).Result()
		if err != nil {
			return nil, err
		}
		if len(values) == 0 {
			return &types.ArticleListRsp{Articles: nil}, nil
		}

		ids = make([]int64, 0, len(values))
		for _, s := range values {
			id, e := strconv.ParseInt(s, 10, 64)
			if e != nil {
				continue
			}
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return &types.ArticleListRsp{Articles: nil}, nil
	}

	// TODO 这里是直接从数据库获取信息，可以修改为先从redis中尝试获取数据之后再去数据库中取
	var dbArticles []model.Article
	if err := l.svcCtx.Db.Where("id IN ?", ids).Find(&dbArticles).Error; err != nil {
		return nil, err
	}
	// 建立 id->article map 以便按 ids 顺序组装响应
	am := make(map[int64]model.Article, len(dbArticles))
	for _, a := range dbArticles {
		am[int64(a.ID)] = a
	}

	rsp := make([]types.ArticleInfo, 0, len(ids))
	for _, id := range ids {
		a, ok := am[id]
		if !ok {
			continue
		}
		var urls model.Pictures
		_ = json.Unmarshal([]byte(a.Pictures), &urls)

		var pic string
		if len(urls.Name) > 0 {
			pic = urls.Name[0]
		}
		info := types.ArticleInfo{
			ID:         a.ID,
			Author:     a.Author,
			Section:    a.Section,
			Time:       a.CreatedAt.Format("2006.01.02 15:04"),
			Summary:    truncateSummary(a.Content, 20),
			Title:      a.Title,
			PictureUrl: pic,
			LookNum:    a.LookNum,
			LikeNum:    a.LikeNum,
			CommentNum: a.CommentNum,
		}
		rsp = append(rsp, info)
	}

	return &types.ArticleListRsp{
		Articles: rsp,
	}, nil
}
