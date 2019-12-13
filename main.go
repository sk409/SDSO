package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func f(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		log.Println(fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
}

func main() {
	v := url.Values{}
	v.Set("name", "aaa")
	v.Set("password", "bbb")
	request, err := http.NewRequest(http.MethodGet, "https://google.com", strings.NewReader(v.Encode()))
	if err != nil {
		panic(err)
	}
	bytes, err := httputil.DumpRequest(request, true)
	fmt.Println(string(bytes))
}
