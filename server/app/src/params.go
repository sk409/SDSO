package main

import (
	"os"
	"path/filepath"

	"github.com/sk409/gogit"
)

var (
	gitClones       *gogit.Git
	gitRepositories *gogit.Git
	gitServer       *gogit.HTTPServer
)

const (
	databaseHost      = "database"
	serverHostAndPort = "0.0.0.0:8080"
	serverScheme      = "http"
	tableNameProjects = "projects"
	tableNameUsers    = "users"
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rootRepositoryPath := filepath.Join(cwd, "..", "repositories")
	gitBinPath := "/usr/bin/git"
	gitClones = gogit.NewGit(filepath.Join(cwd, "..", "clones"), gitBinPath)
	gitRepositories = gogit.NewGit(rootRepositoryPath, gitBinPath)
	gitServer = gogit.NewHTTPServer(rootRepositoryPath, gitBinPath)
}
