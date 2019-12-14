package main

import "net/http"

const (
	cookieNameSessionID = "sessionID"
	cookie30Days        = 60 * 60 * 24 * 30
)

func newCookie(name, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: cookie30Days,
	}
}
