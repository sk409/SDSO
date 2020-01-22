package main

import (
	"reflect"

	"github.com/sk409/gocase"
	"github.com/sk409/gotype"
)

type facade interface {
	public() interface{}
}

func (t test) public() interface{} {
	m, err := convert(t)
	if err != nil {
		return t
	}
	results := []testResult{}
	_, err = find(map[string]interface{}{"test_id": t.ID}, &results)
	if err != nil {
		return t
	}
	rp := make([]interface{}, len(results))
	for index, result := range results {
		var i interface{} = result
		rp[index] = i.(facade).public()
	}
	m["results"] = rp
	return m
}

func (t testResult) public() interface{} {
	m, err := convert(t)
	if err != nil {
		return t
	}
	ts := testStatus{}
	_, err = first(map[string]interface{}{"id": t.TestStatusID}, &ts)
	if err != nil {
		return t
	}
	m["status"] = ts
	return m
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
			l := string(gocase.LowerCamelCase([]byte(key.String()), true))
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
	return data, nil
}
