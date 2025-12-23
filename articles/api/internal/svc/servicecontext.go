// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"errors"
	"fmt"
	"fungo/articles/api/internal/config"
	"fungo/articles/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	Db          *gorm.DB
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化mysql
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("连接到数据库成功")
	}
	_ = db.AutoMigrate(model.Article{})
	_ = db.AutoMigrate(model.Comment{})

	// 初始化 redis
	cli := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
	})
	pong, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("连接redis成功：", pong)

	svc := &ServiceContext{
		Config:      c,
		Db:          db,
		RedisClient: cli,
	}

	// 初始化
	initialize(cli, db)

	// 创建redis到mysql的聚合器
	lookNumAggregator := NewLookNumAggregator(svc, "article-look-nums", 10*time.Second)
	lookNumAggregator.Start()

	likeAggregator := NewLikeAggregator(svc, 10*time.Second)
	likeAggregator.Start()

	return svc
}

// 初始化redis，将数据库中的信息放入redis中
func initialize(cli *redis.Client, db *gorm.DB) {
	cli.FlushAll(context.Background())
	ctx := context.Background()
	sections := []int{2, 3}

	for _, sec := range sections {
		secKey := "section-" + strconv.Itoa(sec) + "-list"
		var arts []model.Article
		if err := db.Where("section = ?", sec).Find(&arts).Error; err != nil {
			fmt.Println("加载 section", sec, "文章失败:", err)
			continue
		}
		z := make([]redis.Z, len(arts))
		for i, _ := range arts {
			score := float64(arts[i].CreatedAt.UnixNano()) / 1e6
			z[i].Score = score
			z[i].Member = arts[i].ID
		}
		cli.ZAdd(ctx, secKey, z...)
		cli.ZAdd(ctx, "section-1-list", z...)
		for i, _ := range arts {
			score := arts[i].LikeNum*2 + arts[i].LookNum + arts[i].CommentNum*5
			z[i].Score = float64(score)
			z[i].Member = arts[i].ID
		}
		cli.ZAdd(ctx, "hot-article-list", z...)
	}
	fmt.Println("初始化redis成功")
}

type LookNumAggregator struct {
	svcCtx   *ServiceContext
	key      string
	interval time.Duration
	stop     chan struct{}
}

// NewLookNumAggregator 创建聚合器，key 为 Redis 哈希键，interval 为聚合周期
func NewLookNumAggregator(svcCtx *ServiceContext, key string, interval time.Duration) *LookNumAggregator {
	return &LookNumAggregator{
		svcCtx:   svcCtx,
		key:      key,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start 启动后台聚合器
func (a *LookNumAggregator) Start() {
	ticker := time.NewTicker(a.interval)
	go func() {
		defer a.Stop()
		for {
			select {
			case <-ticker.C:
				if err := a.aggregateOnce(context.Background()); err != nil {
					fmt.Println("looknum aggregate error:", err)
				}
			case <-a.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop 停止聚合器
func (a *LookNumAggregator) Stop() {
	close(a.stop)
}

// Lua 脚本：原子地获取哈希所有字段并删除该哈希
var hgetallAndDel = `
	local res = redis.call('HGETALL', KEYS[1])
	if next(res) == nil then
		return res
	end
	redis.call('DEL', KEYS[1])
	return res
`

// aggregateOnce 执行一次聚合：从 Redis 读取并清空增量，然后写入 MySQL（批量 CASE WHEN）
func (a *LookNumAggregator) aggregateOnce(ctx context.Context) error {
	r := a.svcCtx.RedisClient
	db := a.svcCtx.Db

	// 原子读取并删除
	raw, err := r.Eval(ctx, hgetallAndDel, []string{a.key}).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	// raw 期望是 []interface{}{field1, val1, field2, val2, ...}
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}

	// 解析到 map[id]delta
	deltas := make(map[uint64]int64, len(arr)/2)
	for i := 0; i+1 < len(arr); i += 2 {
		fb, _ := arr[i].(string)
		vb, _ := arr[i+1].(string)
		if fb == "" || vb == "" {
			continue
		}
		id, err := strconv.ParseUint(fb, 10, 64)
		if err != nil {
			continue
		}
		delta, err := strconv.ParseInt(vb, 10, 64)
		if err != nil {
			continue
		}
		if delta == 0 {
			continue
		}
		deltas[id] += delta
	}

	if len(deltas) == 0 {
		return nil
	}

	for i, v := range deltas {
		if err = db.Model(&model.Article{}).Where("id = ?", i).
			UpdateColumn("look_num", gorm.Expr("look_num + ?", v)).Error; err != nil {
			fmt.Println("db increment look_num error:", err)
			if err2 := r.HSet(ctx, a.key, i, v).Err(); err2 != nil {
				// 如果连回写也失败，记录错误并返回事务错误（上层可报警）
				fmt.Println("looknum restore to redis failed:", err2)
			}
			return err
		}
	}
	return nil
}

// LikeAggregator
type LikeAggregator struct {
	svcCtx   *ServiceContext
	interval time.Duration
	stop     chan struct{}
}

func NewLikeAggregator(svcCtx *ServiceContext, interval time.Duration) *LikeAggregator {
	return &LikeAggregator{
		svcCtx:   svcCtx,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (l *LikeAggregator) Start() {
	ticker := time.NewTicker(l.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := l.aggregateOnce(context.Background()); err != nil {
					fmt.Println("looknum aggregate error:", err)
				}
			case <-l.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (l *LikeAggregator) Stop() {
	close(l.stop)
}

func (l *LikeAggregator) aggregateOnce(ctx context.Context) error {
	// 将点赞表保存到数据库中
	list, err := l.svcCtx.RedisClient.SMembers(ctx, "likebitmap-save-list").Result()
	l.svcCtx.RedisClient.Del(ctx, "likebitmap-save-list")
	if err != nil {
		return err
	}
	for _, v := range list {
		s, err := l.svcCtx.RedisClient.Get(ctx, "likebitmap-"+v).Result()
		if err != nil {
			fmt.Println("获取点赞bitmap失败,err = ", err)
			continue
		}
		if err := l.svcCtx.Db.Model(&model.Article{}).Where("id = ?", v).UpdateColumn("likes_bitmap", s).Error; err != nil {
			fmt.Println("保存点赞bitmap失败,err = ", err)
			continue
		}
	}

	// 保存点赞数到数据库中
	raw, err := l.svcCtx.RedisClient.Eval(ctx, hgetallAndDel, []string{"like-save-set"}).Result()
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}
	deltas := make(map[uint64]int64, len(arr)/2)
	for i := 0; i+1 < len(arr); i += 2 {
		fb, _ := arr[i].(string)
		vb, _ := arr[i+1].(string)
		if fb == "" || vb == "" {
			continue
		}
		id, err := strconv.ParseUint(fb, 10, 64)
		if err != nil {
			continue
		}
		delta, err := strconv.ParseInt(vb, 10, 64)
		if err != nil {
			continue
		}
		if delta == 0 {
			continue
		}
		deltas[id] += delta
	}

	if len(deltas) == 0 {
		return nil
	}

	for i, v := range deltas {
		if err = l.svcCtx.Db.Model(&model.Article{}).Where("id = ?", i).
			UpdateColumn("like_num", gorm.Expr("like_num + ?", v)).Error; err != nil {
			fmt.Println("db increment look_num error:", err)
		}
	}

	return nil
}
