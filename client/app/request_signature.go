package main

import (
	"net/http"
	"sort"
)

var requestSignatures = make([]string, 0)

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
