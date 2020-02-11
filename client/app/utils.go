package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"

	"github.com/sk409/gotype"
)

func closePanic(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		panic(err)
	}
}

func emptyAny(values ...interface{}) bool {
	for _, value := range values {
		if gotype.IsString(value) {
			if value.(string) == "" {
				return true
			}
		} else if gotype.IsSlice(value) {
			rv := reflect.ValueOf(value)
			if rv.Len() == 0 {
				return true
			}
		}
	}
	return false
}

func command(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

func readJSON(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer closePanic(file)
	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, data)
	return err
}

func saveJSON(path string, data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer closePanic(file)
	_, err = file.Write(jsonBytes)
	return err
}
