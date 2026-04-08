package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:tag name"`
	Articles []Article `gorm:"many2many:article_tags;" json:"-"`
}
