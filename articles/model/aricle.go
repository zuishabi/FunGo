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
