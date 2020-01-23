package main

import (
	"time"
)

type branchProtectionRule struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	BranchName string     `gorm:"type:varchar(128);not null"`
	ProjectID  uint       `gorm:"not null"`
}

type build struct {
	Docker []docker
	Steps  []interface{}
}

type commit struct {
	Branchname string
	Date       string
	Diff       string
	Message    string
	SHA1       string
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
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Name      string     `gorm:"type:varchar(128);not null"`
	UserID    uint       `gorm:"not null"`
}

type request struct {
	Text string `gorm:"type:text;not null"`
}

type scan struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	CommitSHA1 string     `gorm:"type:char(40);not null"`
	UserID     uint       `gorm:"not null"`
	ProjectID  uint       `gorm:"not null"`
}

type test struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Steps      int        `gorm:"not null"`
	Branchname string     `gorm:"type:varchar(256); not null;"`
	CommitSHA1 string     `gorm:"type:char(40);not null;unique"`
	ProjectID  uint       `gorm:"not null"`
}

type testStatus struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Text      string     `gorm:"type:varchar(7);unique"`
}

type testResult struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
	Command      string     `gorm:"type:text;not null"`
	Output       string     `gorm:"type:text;"`
	TestID       uint       `gorm:"not null"`
	TestStatusID uint       `gorm:"not null"`
	CompletedAt  *time.Time
}

type user struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
	Name             string     `gorm:"type:varchar(32);not null;unique"`
	Password         string     `gorm:"type:char(60);not null;"`
	ProfileImagePath *string    `gorm:"type:varchar(256);unique"`
}

type vulnerability struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Name        string     `gorm:"type:varchar(32);not null"`
	Description string     `gorm:"type:varchar(128);not null"`
	Path        string     `gorm:"type:varchar(256);not null"`
	Method      string     `gorm:"type:varchar(8);not null"`
	Request     string     `gorm:"type:text;not null"`
	Response    string     `gorm:"type:text;not null"`
	ProjectID   uint       `gorm:"not null"`
	ScanID      uint       `gorm:"not null"`
}
