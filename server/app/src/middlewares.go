package main

import (
	"crypto/sha512"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// func setCORSHeaders(w http.ResponseWriter) {
// 	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 	// w.Header().Set("Access-Control-Allow-Methods", "*")
// 	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 	// w.Header().Set("Access-Control-Allow-Credentials", "true")
// }

// func preflightRequestMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodOptions {
// 			setCORSHeaders(w)
// 		} else {
// 			next.ServeHTTP(w, r)
// 		}
// 	})
// }

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Originをenvから取得
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		next.ServeHTTP(w, r)
	})
}

func gitBasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathParams := mux.Vars(r)
		userName := pathParams["user"]
		user := user{}
		db.Where("name = ?", userName).First(&user)
		if user.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		n, p, ok := r.BasicAuth()
		hashedPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(p)))
		if !ok || user.Name != n || user.Password != hashedPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password."`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			w.Header().Set("Content-Type", "text/plain")
			return
		}
		next.ServeHTTP(w, r)
	})
}
