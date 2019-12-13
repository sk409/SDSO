package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sk409/goconst"
)

func sendRequest(method string, urlString string, params map[string]string) (*http.Response, error) {
	var body io.Reader
	if method == http.MethodGet {
		urlString += "?"
		for key, value := range params {
			urlString += fmt.Sprintf("%s=%s&", key, value)
		}
	} else if method == http.MethodPost {
		values := url.Values{}
		for key, value := range params {
			values.Set(key, value)
		}
		body = strings.NewReader(values.Encode())
	} else {
		log.Fatal("未対応のメソッド")
	}
	// log.Println(urlString)
	request, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_URLENCODED)
	storeScanClient := new(http.Client)
	response, err := storeScanClient.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func sendRequestExist(method string, urlString string, params map[string]string) (bool, error) {
	response, err := sendRequest(
		method,
		urlString,
		params,
	)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	existProjectJSONBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	var existProjectJSON map[string]bool
	err = json.Unmarshal(existProjectJSONBytes, &existProjectJSON)
	if err != nil {
		return false, err
	}
	return existProjectJSON["exist"], nil
}
