package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"size:32;index"`
	Password string `gorm:"size:64"`
}
