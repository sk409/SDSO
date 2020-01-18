package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func loadUser() (*user, error) {
	jsonFilePath := makeFilePath(filepath.Join("auth", "user"))
	jsonBytes, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}
	user := user{}
	err = json.Unmarshal(jsonBytes, &user)
	if err != nil {
		return nil, err
	}
	response, err := sendRequest(
		http.MethodGet,
		route(pathUsers),
		map[string]string{
			"name": user.Name,
		},
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var responseJSON map[string]bool
	err = json.Unmarshal(responseBytes, &responseJSON)
	if err != nil {
		return nil, err
	}
	if len(responseJSON) == 0 {
		return nil, errors.New("The specified user does not exist")
	}
	return &user, nil
}
