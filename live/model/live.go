package model

type UserRoom struct {
	UID    uint64 `gorm:"primaryKey"`
	RoomID uint64 `gorm:"UniqueIndex"`
}

type RoomInfo struct {
	RoomID      uint64 `gorm:"primaryKey"`
	Title       string
	Cover       string
	Description string
}
