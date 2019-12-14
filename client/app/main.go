package main

import (
	"os"
)

func init() {
	mkDirIfNotExist(makeFilePath(""))
}

func main() {
	if len(os.Args) == 1 {
		return
	}
	subcommand := os.Args[1]
	switch subcommand {
	case "login":
		entrypointLogin()
	case "record":
		entrypointRecord()
	case "request":
		entrypointRequest()
	case "scan":
		entrypointScan()
	case "upload":
		entrypointUpload()
	case "x509":
		entrypointX509()
	}
}
