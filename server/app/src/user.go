package main

import (
	"github.com/jinzhu/gorm"
)

type user struct {
	gorm.Model
	Name     string `gorm:"type:varchar(32);not null"`
	Password string `gorm:"type:varchar(512);not null"`
}
