package main

import "github.com/jinzhu/gorm"

type project struct {
	gorm.Model
	Name   string `gorm:"type:varchar(128);not null"`
	UserID uint   `gorm:"not null"`
}

type request struct {
	Text string `gorm:"type:text;not null"`
}

type scan struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	ProjectID uint `gorm:"not null"`
}

type user struct {
	gorm.Model
	Name     string `gorm:"type:varchar(32);not null"`
	Password string `gorm:"type:varchar(512);not null"`
}

type vulnerability struct {
	gorm.Model
	Name        string `gorm:"type:varchar(32);not null"`
	Description string `gorm:"type:varchar(128);not null"`
	Path        string `gorm:"type:varchar(256);not null"`
	Method      string `gorm:"type:varchar(8);not null"`
	Request     string `gorm:"type:text;not null"`
	Response    string `gorm:"type:text;not null"`
	ProjectID   uint   `gorm:"not null"`
	ScanID      uint   `gorm:"not null"`
}
