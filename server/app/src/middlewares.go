package main

import (
	"net/http"
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
