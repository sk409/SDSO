package main

import (
	"os"
	"path/filepath"
)

func existDirectory(path string) bool {
	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		return false
	}
	return true
}

func makeFilePath(path string) string {
	return filepath.Join("/etc", "sdso", path)
}

func mkDirIfNotExist(path string) bool {
	if existDirectory(path) {
		return false
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return false
	}
	return true
}
