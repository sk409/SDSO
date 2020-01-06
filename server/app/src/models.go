package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type branchProtectionRule struct {
	gorm.Model
	BranchName string `gorm:"type:varchar(128);not null"`
	ProjectID  uint   `gorm:"not null"`
}

type build struct {
	Docker []docker
	Steps  []interface{}
}

type config struct {
	Version int
	Jobs    jobs
}

type docker struct {
	Image string
}

type jobs struct {
	Build build
}

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

type test struct {
	gorm.Model
	Steps     int  `gorm:"not null"`
	ProjectID uint `gorm:"not null"`
}

type testStatus struct {
	gorm.Model
	Text string `gorm:"type:varchar(7);unique"`
}

type testResult struct {
	gorm.Model
	Command      string `gorm:"type:text;not null"`
	Output       string `gorm:"type:text;"`
	TestID       uint   `gorm:"not null"`
	TestStatusID uint   `gorm:"not null"`
	CompletedAt  *time.Time
}

type user struct {
	gorm.Model
	Name             string  `gorm:"type:varchar(32);not null;unique"`
	Password         string  `gorm:"type:varchar(512);not null;"`
	ProfileImagePath *string `gorm:"type:varchar(256);unique"`
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
