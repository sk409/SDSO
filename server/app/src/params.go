package main

import (
	"os"
	"path/filepath"

	"github.com/sk409/gogit"
)

var (
	cwd             string
	gitClones       *gogit.Git
	gitRepositories *gogit.Git
	gitTesting      *gogit.Git
	gitServer       *gogit.HTTPServer
)

const (
	databaseHost      = "database"
	gitBinPath        = "/usr/bin/git"
	serverHostAndPort = "0.0.0.0:8080"
	serverScheme      = "http"
	tableNameProjects = "projects"
	tableNameUsers    = "users"
)

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	rootRepositoryPath := filepath.Join(cwd, "..", "repositories")
	gitClones = gogit.NewGit(filepath.Join(cwd, "..", "clones"), gitBinPath)
	gitRepositories = gogit.NewGit(rootRepositoryPath, gitBinPath)
	//gitTesting = gogit.NewGit(filepath.Join(cwd, "..", "testing"))
	gitServer = gogit.NewHTTPServer(rootRepositoryPath, gitBinPath)
}
