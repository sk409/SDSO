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
	userJSONFilePath         = ""
)

const (
	directoryProject    = ".sdso"
	pathProjects        = "/projects"
	pathScans           = "/scans"
	pathUsers           = "/users"
	pathVulnerabilities = "/vulnerabilities"
	serverOrigin        = "http://server:8080"
)

func init() {
	directoryApp = filepath.Join("/etc", "sdso")
	directoryCA = filepath.Join(directoryApp, "ca")
	directoryAuth = filepath.Join(directoryApp, "auth")
	directoryRequests = filepath.Join(directoryProject, "requests")
	directoryVulnerabilities = filepath.Join(directoryProject, "vulnerabilities")
	filepathConfig = filepath.Join(directoryProject, "config.json")
	filepathUser = filepath.Join(directoryAuth, "user")
	userJSONFilePath = filepath.Join(directoryAuth, "user")
}
