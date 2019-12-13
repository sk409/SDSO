package main

import (
	"github.com/jinzhu/gorm"
)

type scan struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	ProjectID uint `gorm:"not null"`
}
