package model

import "time"

type Logs struct {
	ID       uint      `gorm:"primaryKey"`
	CreateAt time.Time `gorm:"autoCreateTime"`
	UpdateAt time.Time `gorm:"autoUpdateTime"`
	Hostname string
	File     string
	Log      string
}

func (Logs) TableName() string {
	return "logs"
}
