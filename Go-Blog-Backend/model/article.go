package model

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	// article
	Title     string `gorm:"type:varchar(128);not null;index;comment:article title"`
	Summary   string `gorm:"type:varchar(255);not null;comment:article summary"`
	Content   string `gorm:"type:longtext;not null;comment:article content"`
	ViewCount uint   `gorm:"default:0;comment:article views"`
	IsDraft   bool   `gorm:"default:true;comment:status(true indicates a draft; false indicates published)"`
	Id        int64  `gorm:"primaryKey;autoIncrement:false" json:"ID,string"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// author
	UserId int64 `gorm:"not null;comment: Primary Key relation with user"`
	User   User  `gorm:"foreignKey:UserID"`

	// category
	CategoryId uint     `gorm:"comment: Primary Key relation with category"`
	Category   Category `gorm:"foreignKey:CategoryID"`

	// tag
	Tags []Tag `gorm:"many2many:article_tags;"`
}
