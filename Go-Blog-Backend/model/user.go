package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username  string `gorm:"type:varchar(50);uniqueIndex;not null;comment:admin account username"`
	Password  string `gorm:"type:varchar(255);not null;comment:admin account password"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null;comment:admin email"`
	Id        int64  `gorm:"primaryKey;autoIncrement:false" json:"ID,string"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
