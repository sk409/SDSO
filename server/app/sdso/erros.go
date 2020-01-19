package main

import "errors"

var (
	errBadRequest  = errors.New("Bad Request")
	errNotExist    = errors.New("Not Exist")
	errInvalidType = errors.New("Invalid Type")
)
