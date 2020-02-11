package main

type config struct {
	Teamname    string `json:"teamname"`
	Projectname string `json:"projectname"`
}

type request struct {
	Method string              `json:"method"`
	URL    string              `json:"url"`
	Path   string              `json:"path"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

type scan struct {
	ID uint `json:"id"`
}

type user struct {
	ID        uint   `json:"ID"`
	Name      string `json:"Name"`
	Password  string `json:"Password"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
}

type vulnerability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	Request     string `json:"request"`
	Response    string `json:"response"`
	ScanID      uint   `json:"scanID"`
}
