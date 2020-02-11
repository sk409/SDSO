package main

import (
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

var requestSignatures = make([]string, 0)

func (r *request) toHTTPRequestWithSignature(signature string) (*http.Request, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}
	if r.Method == http.MethodGet {
		query, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			return nil, err
		}
		for key, values := range query {
			query.Del(key)
			for range values {
				query.Add(key, signature)
			}
		}
		u.RawQuery = query.Encode()
	}
	var body io.Reader
	if r.Method == http.MethodPost {
		// TODO: Content-Typeによって場合わけ
		query, err := url.ParseQuery(r.Body)
		if err != nil {
			return nil, err
		}
		urlValues := url.Values{}
		for key, values := range query {
			for range values {
				urlValues.Add(key, signature)
			}
		}
		body = strings.NewReader(urlValues.Encode())
	}
	request, err := http.NewRequest(r.Method, u.String(), body)
	for key, values := range r.Header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	if err != nil {
		return nil, err
	}
	return request, nil
}

func newRequestSignature(r *http.Request) string {
	requestSignature := r.URL.Path
	if r.Method == http.MethodGet {
		keys := []string{}
		for key := range r.URL.Query() {
			keys = append(keys, key)
		}
		sort.Sort(sort.StringSlice(keys))
		for _, key := range keys {
			requestSignature += key
		}
	}
	return requestSignature
}

func storeRequestSignature(requestSignature string) {
	requestSignatures = append(requestSignatures, requestSignature)
}

func isNewRequest(r *http.Request) (string, bool) {
	requestSignature := newRequestSignature(r)
	for _, rs := range requestSignatures {
		if rs == requestSignature {
			return "", false
		}
	}
	return requestSignature, true
}
