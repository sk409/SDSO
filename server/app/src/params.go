package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"

	"github.com/sk409/gogit"
)

var (
	cwd                  string
	gitClones            *gogit.Git
	gitRepositories      *gogit.Git
	gitTesting           *gogit.Git
	gitServer            *gogit.HTTPServer
	websocketsTest       = map[uint]*websocket.Conn{}
	websocketsTestResult = map[uint]*websocket.Conn{}
	websocketUpgrader    = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	databaseHost      = "database"
	gitBinPath        = "/usr/bin/git"
	messageBufferSize = 256
	serverHostAndPort = "0.0.0.0:8080"
	serverScheme      = "http"
	socketBufferSize  = 1024
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
