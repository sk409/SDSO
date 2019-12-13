package main

import (
	"os"
	"path/filepath"

	"github.com/sk409/gogit"
)

var (
	gitServer *gogit.HTTPServer
)

const (
	databaseHost      = "database"
	tableNameUsers    = "users"
	tableNameProjects = "projects"
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	gitServer = gogit.NewHTTPServer(filepath.Join(cwd, "..", "repositories"), "/usr/bin/git")
}
