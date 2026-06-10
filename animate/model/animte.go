package model

import "time"

type AnimateList struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Description string
	Year        string
	Tags        string
	State       string
	Num         int
}

type TodayUpdateList struct {
	ID        uint64 `gorm:"primaryKey"`
	UpdatedAt time.Time
}

type SubscribeFavoriteList struct {
	UID       uint64 `gorm:"primaryKey;autoIncrement:false"`
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	Subscribe bool
	Favorite  bool
}

type AnimateUpdateInfo struct {
	ID          uint64 `gorm:"primaryKey"`
	Version     int
	Name        string
	Description string
	UpdatedAt   time.Time
}

type WishList struct {
	ID        uint64 `gorm:"primaryKey"`
	Content   string
	UID       uint64 `gorm:"index"`
	CreatedAt time.Time
}
