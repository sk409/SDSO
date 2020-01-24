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
	gitTmpRepositories   *gogit.Git
	gitRepositories      *gogit.Git
	gitTesting           *gogit.Git
	gitServer            *gogit.HTTPServer
	testResultColors     = map[string]string{}
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
	databaseHost           = "database"
	gitBinPath             = "/usr/bin/git"
	messageBufferSize      = 256
	serverHostAndPort      = "0.0.0.0:8080"
	serverScheme           = "http"
	serverOrigin           = serverScheme + "://" + serverHostAndPort
	socketBufferSize       = 1024
	tableNameProjects      = "projects"
	tableNameUsers         = "users"
	testResultFailedColor  = "rgb(220, 102, 97)"
	testResultFailedText   = "failed"
	testResultRunningColor = "rgb(130, 209, 226)"
	testResultRunningText  = "running"
	testResultSuccessColor = "rgb(107, 197, 143)"
	testResultSuccessText  = "success"
)

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	rootRepositoryPath := filepath.Join(cwd, "..", "repositories")
	gitTmpRepositories = gogit.NewGit(filepath.Join(cwd, "..", "tmp_repositories"), gitBinPath)
	gitRepositories = gogit.NewGit(filepath.Join(cwd, "..", "repositories"), gitBinPath)
	gitServer = gogit.NewHTTPServer(rootRepositoryPath, gitBinPath)
	testResultColors[testResultFailedText] = testResultFailedColor
	testResultColors[testResultRunningText] = testResultRunningColor
	testResultColors[testResultSuccessText] = testResultSuccessColor
}
