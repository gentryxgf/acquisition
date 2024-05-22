package model

import "time"

type UsedPercent struct {
	ID        uint      `gorm:"primaryKey"`
	CreateAt  time.Time `gorm:"autoCreateTime"`
	UpdateAt  time.Time `gorm:"autoUpdateTime"`
	Metric    string
	Endpoint  string
	Timestamp int
	Step      int
	Value     float32
	Extend    string
}

func (UsedPercent) TableName() string {
	return "used_percent"
}
