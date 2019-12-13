package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type request struct {
	Method string              `json:"method"`
	URL    string              `json:"url"`
	Path   string              `json:"path"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

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
