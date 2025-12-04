package model

import "time"

type Article struct {
	ID          uint64 `gorm:"primaryKey"`
	Content     string
	CreatedAt   time.Time
	Section     uint32 `gorm:"Index"`
	Author      uint64 `gorm:"Index"`
	Pictures    string
	Title       string
	LookNum     int
	LikeNum     int
	LikesBitmap []byte `gorm:"column:likes_bitmap;type:blob"`
	CommentNum  int
}

type Pictures struct {
	Name []string `json:"name"`
}

type Comment struct {
	ID        uint64 `gorm:"primaryKey"`
	UID       uint64 `gorm:"Index"`
	ArticleID uint64 `gorm:"Index"`
	Parent    uint64
	CreatedAt time.Time
	PParent   uint64 // 在回复评论中回复
	Content   string
}
