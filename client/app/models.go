package main

type config struct {
	Teamname    string `json:"teamname"`
	Projectname string `json:"projectname"`
}

type scan struct {
	ID uint `json:"id"`
}
