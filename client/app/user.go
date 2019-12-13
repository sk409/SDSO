package main

type user struct {
	ID        uint   `json:"ID"`
	Name      string `json:"Name"`
	Password  string `json:"Password"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
}
