package main

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func hookRequest(r *http.Request) {
	if len(targetHost) == 0 || !strings.Contains(r.Host, targetHost) {
		return
	}
	requestSignature, isNewRequest := isNewRequest(r)
	if !isNewRequest {
		return
	}
	storeRequestSignature(requestSignature)
	header := make(map[string][]string)
	for key, values := range r.Header {
		for _, value := range values {
			header[key] = append(header[key], value)
		}
	}
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		return
	}
	requestString := string(requestBytes)
	headerAndBody := strings.Split(requestString, "\r\n\r\n")
	bodyString := ""
	if len(headerAndBody) == 2 {
		bodyString = headerAndBody[1]
	}
	request := request{
		Method: r.Method,
		URL:    r.URL.String(),
		Path:   r.URL.Path,
		Header: header,
		Body:   bodyString,
	}
	jsonBytes, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		return
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}
	file, err := os.Create(filepath.Join("requests", id.String()))
	if err != nil {
		return
	}
	defer file.Close()
	file.Write(jsonBytes)
}

func hookResponse(r *http.Response) {
	// if len(targetHost) == 0 || !strings.Contains(r.Request.Host, targetHost) {
	// 	return
	// }
	// bytes, err := httputil.DumpResponse(r, true)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(string(bytes))
}
