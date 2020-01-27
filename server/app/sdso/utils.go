package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"

	"github.com/sk409/gocase"
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

func find(query map[string]interface{}, model interface{}) error {
	db.Where(query).Find(model)
	return db.Error
}

func findByUniqueKey(uniqueKeys interface{}, model interface{}) error {
	db.Where(uniqueKeys).Find(model)
	return db.Error
}

func first(query map[string]interface{}, model interface{}) error {
	db.Where(query).First(model)
	return db.Error
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

func save(query map[string]interface{}, model interface{}) error {
	rv := reflect.ValueOf(model).Elem()
	for key, value := range query {
		fieldName := string(gocase.UpperCamelCase([]byte(key), true))
		fv := rv.FieldByName(fieldName)
		ft := fv.Type()
		// vt := reflect.TypeOf(value)
		log.Println(fieldName)
		log.Println(reflect.TypeOf(value))
		log.Println("======")
		if ft.Kind() == reflect.String {
			fv.SetString(value.(string))
		} else if ft.Kind() == reflect.Uint {
			var v interface{}
			if gotype.IsString(value) {
				s := value.(string)
				var err error
				v, err = strconv.ParseUint(s, 10, 64)
				if err != nil {
					return err
				}
			} else {
				v = value
			}
			fv.SetUint(v.(uint64))
		}
	}
	db.Save(model)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
