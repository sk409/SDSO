package main

import (
	"os"
)

func main() {
	if len(os.Args) == 1 {
		return
	}
	subcommand := os.Args[1]
	switch subcommand {
	case "init":
		entrypointInit()
	case "login":
		entrypointLogin()
	case "push":
		entrypointPush()
	case "record":
		entrypointRecord()
	case "setup":
		entrypointSetup()
	case "scan":
		entrypointScan()
	case "upload":
		entrypointUpload()
	case "x509":
		entrypointX509()
	}
}
