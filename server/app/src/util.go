package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/sk409/goconst"
)

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

func getBranchName(r *http.Request) (string, error) {
	var body io.ReadCloser
	var err error
	if r.Header.Get(goconst.HTTP_HEADER_CONTENT_ENCODING) == goconst.HTTP_HEADER_CONTENT_ENCODING_GZIP {
		body, err = gzip.NewReader(r.Body)
		if err != nil {
			return "", err
		}
		defer body.Close()
	} else {
		body = r.Body
	}
	requestBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	branchNameRegex := regexp.MustCompile("refs/heads/([^ ]+)")
	branchNameSubmatches := branchNameRegex.FindSubmatch(requestBytes)
	if len(branchNameSubmatches) < 2 {
		return "", nil
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBytes))
	branchName := string(branchNameSubmatches[1])
	//
	// TODO: なぜ?
	branchName = branchName[:len(branchName)-1]
	//
	return branchName, nil
}
