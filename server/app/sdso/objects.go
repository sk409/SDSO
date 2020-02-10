package main

type build struct {
	Docker []docker
	Steps  []interface{}
}

type commit struct {
	Branchname string `json:"commit"`
	Date       string `json:"date"`
	Diff       string `json:"diff"`
	Message    string `json:"message"`
	SHA1       string `json:"sha1"`
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
