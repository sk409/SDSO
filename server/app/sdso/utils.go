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

	"github.com/sk409/gocase"
	"github.com/sk409/goconst"
	"github.com/sk409/gotype"
)

type printable interface {
	String() string
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

// func find(query map[string]interface{}, model interface{}) error {
// 	q := map[string]interface{}{}
// 	for key, value := range query {
// 		s := string(gocase.SnakeCase([]byte(key)))
// 		q[s] = value
// 	}
// 	gormDB.Where(q).Find(model)
// 	return gormDB.Error
// }

// func findByUniqueKey(uniqueKeys interface{}, model interface{}) error {
// 	gormDB.Where(uniqueKeys).Find(model)
// 	return gormDB.Error
// }

// func first(query map[string]interface{}, model interface{}) error {
// 	gormDB.Where(query).First(model)
// 	return gormDB.Error
// }

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
		fieldname := string(gocase.UpperCamelCase([]byte(key), true))
		fv := rv.FieldByName(fieldname)
		ft := fv.Type()
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.String {
			fv.SetString(value.(string))
		} else if ft.Kind() == reflect.Int {
			if gotype.IsInt(value) {
				fv.SetInt(int64(value.(int)))
			}
		} else if ft.Kind() == reflect.Uint {
			var v uint
			if gotype.IsString(value) {
				s := value.(string)
				ui, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return err
				}
				v = uint(ui)
			} else if gotype.IsUint(value) {
				v = value.(uint)
			} else {
				continue
			}
			if fv.Kind() == reflect.Ptr {
				fv.Set(reflect.ValueOf(&v))
			} else {
				fv.SetUint(uint64(v))
			}
		}
	}
	gormDB.Save(model)
	if gormDB.Error != nil {
		return gormDB.Error
	}
	return nil
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
