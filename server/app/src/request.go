package main

type request struct {
	Text string `gorm:"type:text;not null"`
}
