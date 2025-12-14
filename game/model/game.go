package model

import "time"

type Links map[string]string

type Pictures struct {
	Name []string `json:"name"`
}
type GameInfo struct {
	ID               uint64 `gorm:"primaryKey"`
	Title            string
	UpdatedAt        time.Time
	Author           uint64
	Description      string
	Links            string
	Pictures         string
	Cover            string
	DownloadFile     string
	CanOnlinePlaying bool
	PlayTime         int
}
