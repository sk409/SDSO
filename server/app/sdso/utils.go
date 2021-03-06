package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/sk409/goconst"
	"github.com/sk409/gotype"
)

type printable interface {
	String() string
}

func checkPermission(u *user, users []user) bool {
	ok := false
	for _, user := range users {
		if u.ID == user.ID {
			ok = true
			break
		}
	}
	return ok
}

func command(directory, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = directory
	return cmd.Output()
}

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
	if rt.Name() == "Time" {
		t := data.(time.Time)
		return t.Format(time.RFC3339Nano), nil
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

func reverse(s interface{}) (interface{}, error) {
	if !gotype.IsSlice(s) {
		return nil, errInvalidType
	}
	rv := reflect.ValueOf(s)
	for i, j := 0, rv.Len()-1; i < j; i, j = i+1, j-1 {
		tmp := reflect.ValueOf(rv.Index(i).Interface())
		rv.Index(i).Set(rv.Index(j))
		rv.Index(j).Set(tmp)
	}
	return s, nil
}

func stringsToUints(strings []string) []uint {
	uints := make([]uint, len(strings))
	for index, str := range strings {
		u, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			continue
		}
		uints[index] = uint(u)
	}
	return uints
}
