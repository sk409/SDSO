package main

import (
	"reflect"

	"github.com/sk409/gocase"
	"github.com/sk409/gotype"
)

type facade interface {
	public() interface{}
}

func public(data interface{}) (interface{}, error) {
	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	if gotype.IsMap(data) {
		if rt.Key().Kind() != reflect.String {
			return nil, errInvalidType
		}
		m := make(map[string]interface{})
		for _, key := range rv.MapKeys() {
			l := string(gocase.LowerCamelCase([]byte(key.String()), false))
			p, err := public(rv.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
			m[l] = p
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
		var i interface{}
		if f, ok := data.(facade); ok {
			i = f.public()
		} else {
			var err error
			i, err = convert(data)
			if err != nil {
				return nil, err
			}
		}
		p, err := public(i)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	return data, nil
}
