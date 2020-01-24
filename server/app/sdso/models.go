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

func (s scan) public() interface{} {
	i, err := convert(s)
	if err != nil {
		return s
	}
	m := i.(map[string]interface{})
	vulnerabilities := []vulnerability{}
	_, err = find(map[string]interface{}{"scan_id": s.ID}, &vulnerabilities)
	if err != nil {
		m["vulnerabilities"] = []vulnerability{}
		return m
	}
	vp := make([]interface{}, len(vulnerabilities))
	for index, v := range vulnerabilities {
		p, err := public(v)
		if err != nil {
			continue
		}
		vp[index] = p
	}
	m["vulnerabilities"] = vp
	return m
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

func (t test) public() interface{} {
	i, err := convert(t)
	if err != nil {
		return t
	}
	m := i.(map[string]interface{})
	results := []testResult{}
	_, err = find(map[string]interface{}{"test_id": t.ID}, &results)
	if err != nil {
		return t
	}
	rp := make([]interface{}, len(results))
	for index, result := range results {
		rp[index] = result.public()
	}
	m["results"] = rp
	return m
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

func (t testResult) public() interface{} {
	i, err := convert(t)
	if err != nil {
		return t
	}
	m := i.(map[string]interface{})
	ts := testStatus{}
	_, err = first(map[string]interface{}{"id": t.TestStatusID}, &ts)
	if err != nil {
		return t
	}
	m["status"] = ts
	return m
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
