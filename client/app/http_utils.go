package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"

	"github.com/sk409/goconst"

	"github.com/sk409/gotype"
)

func fetch(p string, query map[string]interface{}, model interface{}) error {
	u := serverOrigin + p
	u += "?"
	for key, value := range query {
		rv := reflect.ValueOf(value)
		if gotype.IsSlice(value) {
			if !strings.HasSuffix(key, "[]") {
				key += "[]"
			}
			for i := 0; i < rv.Len(); i++ {
				ev := rv.Index(i)
				u += fmt.Sprintf("%s=%s", key, ev.String())
				u += "&"
			}
		} else {
			u += fmt.Sprintf("%s=%s", key, rv.String())
			u += "&"
		}
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_URLENCODED)
	if err != nil {
		return err
	}
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closePanic(res.Body)
	jsonBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, model)
	return err
}

func store(p string, data map[string]interface{}, model interface{}) error {
	values := url.Values{}
	for key, value := range data {
		rv := reflect.ValueOf(value)
		if gotype.IsSlice(value) {
			if !strings.HasSuffix(key, "[]") {
				key += "[]"
			}
			for i := 0; i < rv.Len(); i++ {
				ev := rv.Index(i)
				values.Add(key, ev.String())
			}
		} else {
			values.Set(key, rv.String())
		}
	}
	u := serverOrigin + path.Join(p)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(values.Encode()))
	req.Header.Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_URLENCODED)
	if err != nil {
		return err
	}
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer closePanic(res.Body)
	jsonBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, model)
	return err
}
