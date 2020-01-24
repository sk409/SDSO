package main

import (
	"path/filepath"
)

var (
	directoryApp             = ""
	directoryAuth            = ""
	directoryCA              = ""
	directoryRequests        = ""
	directoryVulnerabilities = ""
	filepathConfig           = ""
	filepathUser             = ""
	targetHost               = ""
)

const (
	pathProjects        = "projects"
	pathScans           = "/scans"
	pathUsers           = "users"
	pathVulnerabilities = "/vulnerabilities"
	projectDirectory    = ".sdso"
	serverOrigin        = "http://server:8080"
	// pathUsersExist    = "users/exist"
	// pathProjectsExist = "projects/exist"
)

func init() {
	directoryApp = filepath.Join("/etc", "sdso")
	directoryCA = filepath.Join(directoryApp, "ca")
	directoryAuth = filepath.Join(directoryApp, "auth")
	directoryRequests = filepath.Join(directoryApp, "requests")
	directoryVulnerabilities = filepath.Join(directoryApp, "vulnerabilities")
	filepathConfig = filepath.Join(projectDirectory, "config.json")
	filepathUser = filepath.Join(directoryAuth, "user")
}
