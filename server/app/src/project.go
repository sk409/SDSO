package main

import (
	"github.com/jinzhu/gorm"
)

type project struct {
	gorm.Model
	Name   string `gorm:"type:varchar(128);not null"`
	UserID uint   `gorm:"not null"`
}
