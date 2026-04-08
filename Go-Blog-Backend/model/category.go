package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);uniqueIndex;not null;comment:category name"`
	Desc    string `gorm:"type:varchar(255);comment:description"`
	Article []Article
}
