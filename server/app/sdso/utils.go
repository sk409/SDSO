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

	"github.com/sk409/goconst"
	"github.com/sk409/gotype"
)

type printable interface {
	String() string
}

// func contains(sequence interface{}, value interface{}) (bool, error) {
// 	if !gotype.IsSlice(sequence) {
// 		return false, errInvalidType
// 	}
// 	st := reflect.TypeOf(sequence)
// 	sv := reflect.ValueOf(sequence)
// 	vt := reflect.TypeOf(value)
// 	if st.Elem().Kind() == vt.Kind() {
// 		return false, errInvalidType
// 	}
// 	for i := 0; i < st.Len(); i++ {
// 		if reflect.DeepEqual(value, sv.Index(i)) {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

func convert(data interface{}) (interface{}, error) {
	if !gotype.IsStruct(data) && !gotype.IsPointer(data) {
		return nil, errInvalidType
	}
	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	if gotype.IsPointer(data) {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	m := make(map[string]interface{})
	for index := 0; index < rv.NumField(); index++ {
		ft := rt.Field(index)
		fv := rv.Field(index)
		if fv.CanInterface() {
			m[ft.Name] = fv.Interface()
		}
	}
	if len(m) != 0 {
		return m, nil
	}
	p, ok := data.(printable)
	if ok {
		return p.String(), nil
	}
	return nil, nil
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
