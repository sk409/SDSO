package main

import "os"

func existDirectory(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !os.IsNotExist(err) && f.IsDir()
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

func isDirectory(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}
