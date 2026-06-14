// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"fmt"
	"fungo/animate/api/internal/config"
	"fungo/animate/model"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
	*AnimateServer
}

type AnimateServer struct {
	TodayUpdate     []uint64
	TodayUpdateLock sync.RWMutex
}

// 初始化服务器
func (a *AnimateServer) init(db *gorm.DB) {
	// 首先从mysql中获取今日更新
	a.checkTodayUpdateOnce(db)
	go a.checkTodayUpdate(db)
}

func (a *AnimateServer) checkTodayUpdateOnce(db *gorm.DB) {
	a.TodayUpdate = make([]uint64, 0)

	var updates []model.TodayUpdateList
	db.Find(&updates)
	for _, v := range updates {
		if v.UpdatedAt.Format("2006.01.02") == time.Now().Format("2006.01.02") {
			a.TodayUpdate = append(a.TodayUpdate, v.ID)
		} else {
			// 将这条数据从数据库中删除
			db.Where("id = ?", v.ID).Delete(&model.TodayUpdateList{})
		}
	}
}

func (a *AnimateServer) checkTodayUpdate(db *gorm.DB) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.checkTodayUpdateOnce(db)
		}
	}
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("连接到数据库成功")
	db.AutoMigrate(&model.TodayUpdateList{})
	db.AutoMigrate(&model.AnimateList{})
	db.AutoMigrate(&model.SubscribeFavoriteList{})
	db.AutoMigrate(&model.AnimateUpdateInfo{})
	db.AutoMigrate(&model.WishList{})

	animateServer := &AnimateServer{}
	animateServer.init(db)

	return &ServiceContext{
		Config:        c,
		Db:            db,
		AnimateServer: animateServer,
	}
}
