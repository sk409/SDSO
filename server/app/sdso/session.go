package main

import (
	"github.com/sk409/gosession"
)

var sessionManager *gosession.Manager

const (
	sessionStoreNameUserID = "userID"
)

func init() {
	sessionManager = &gosession.Manager{Provider: gosession.NewMemoryProvider()}
}
