package main

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
