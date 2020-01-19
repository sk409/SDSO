package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"

	"github.com/sk409/gocase"

	"github.com/sk409/goconst"
	"github.com/sk409/gotype"
)

func convert(data interface{}) (map[string]interface{}, error) {
	if !gotype.IsStruct(data) && !gotype.IsPointer(data) {
		return nil, errInvalidType
	}
	m := make(map[string]interface{})
	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	if gotype.IsPointer(data) {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	for index := 0; index < rv.NumField(); index++ {
		ft := rt.Field(index)
		fv := rv.Field(index)
		if fv.CanInterface() {
			m[ft.Name] = fv.Interface()
		}
	}
	return m, nil
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

func getBranchNameAndCommitSHA1(r *http.Request) (string, string, error) {
	var body io.ReadCloser
	var err error
	if r.Header.Get(goconst.HTTP_HEADER_CONTENT_ENCODING) == goconst.HTTP_HEADER_CONTENT_ENCODING_GZIP {
		body, err = gzip.NewReader(r.Body)
		if err != nil {
			return "", "", err
		}
		defer body.Close()
	} else {
		body = r.Body
	}
	requestBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "", "", err
	}
	branchNameRegex := regexp.MustCompile("[0-9a-z]+ ([0-9a-z]+) refs/heads/([^ ]+)")
	branchNameSubmatches := branchNameRegex.FindSubmatch(requestBytes)
	if len(branchNameSubmatches) != 3 {
		return "", "", errors.New("Bad request")
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBytes))
	branchName := string(branchNameSubmatches[2])
	commitSHA1 := string(branchNameSubmatches[1])
	//
	// TODO: なぜ?
	branchName = branchName[:len(branchName)-1]
	//
	return branchName, commitSHA1, nil
}

func public(data interface{}) (interface{}, error) {
	if !gotype.IsMap(data) && !gotype.IsSlice(data) && !gotype.IsStruct(data) {

	}
	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	if gotype.IsMap(data) {
		if rt.Key().Kind() != reflect.String {
			return nil, errInvalidType
		}
		m := make(map[string]interface{})
		for _, key := range rv.MapKeys() {
			l := string(gocase.LowerCamelCase([]byte(key.String()), true))
			m[l] = rv.MapIndex(key).Interface()
		}
		return m, nil
	}
	if gotype.IsSlice(data) {
		s := make([]interface{}, rv.Len())
		for index := 0; index < rv.Len(); index++ {
			v := rv.Index(index)
			p, err := public(v.Interface())
			if err != nil {
				return nil, errInvalidType
			}
			s[index] = p
		}
		return s, nil
	}
	if gotype.IsStruct(data) {
		c, err := convert(data)
		if err != nil {
			return nil, err
		}
		p, err := public(c)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	return nil, errInvalidType
	// s := make(map[string]interface{})
	// if !gotype.IsMap(data) {
	// 	var err error
	// 	s, err = convert(data)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// m := make(map[string]interface{})
	// for key, value := range s {
	// 	l := string(gocase.LowerCamelCase([]byte(key), true))
	// 	m[l] = value
	// }
	// return m, nil
}
