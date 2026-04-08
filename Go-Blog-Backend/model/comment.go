package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int64 `gorm:"primaryKey;autoIncrement:false" json:"ID,string"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Nickname  string `json:"nickname" gorm:"type:varchar(20);not null"`
	Content   string `json:"content" gorm:"type:text;not null"`
	IP        string `json:"ip" gorm:"type:varchar(50);index"`
}

type IPStats struct {
	gorm.Model
	IP           string `gorm:"primaryKey"`
	CommentCount int    `gorm:"default:0"`
	IsForbidden  bool   `gorm:"default:false"`
}
