package main

// import (
// 	"net/http"
// 	"path"
// 	"strings"

// 	"golang.org/x/crypto/bcrypt"

// 	"github.com/sk409/goconst"
// )

// // func setCORSHeaders(w http.ResponseWriter) {
// // 	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// // 	// w.Header().Set("Access-Control-Allow-Methods", "*")
// // 	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// // 	// w.Header().Set("Access-Control-Allow-Credentials", "true")
// // }

// // func preflightRequestMiddleware(next http.Handler) http.Handler {
// // 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// // 		if r.Method == http.MethodOptions {
// // 			setCORSHeaders(w)
// // 		} else {
// // 			next.ServeHTTP(w, r)
// // 		}
// // 	})
// // }

// func allowCredentials(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
// 		next.ServeHTTP(w, r)
// 	})
// }

// func allowHeaders(headers []string, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		v := strings.Join(headers, ",")
// 		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, v)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func allowMethods(methods []string, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		v := strings.Join(methods, ",")
// 		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, v)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func cors(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		//TODO: Originをenvから取得
// 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		next.ServeHTTP(w, r)
// 	})
// }

// func gitBasicAuth(next http.Handler) http.Handler {
// 	Unauthorized := func(w http.ResponseWriter) {
// 		w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password."`)
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
// 		w.Header().Set("Content-Type", "text/plain")
// 	}
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		components := strings.Split(r.URL.Path, "/")
// 		if len(components) < 3 {
// 			respond(w, http.StatusNotFound)
// 			return
// 		}
// 		teamname := components[1]
// 		projectname := components[2]
// 		if r.URL.Path == "/"+path.Join(teamname, projectname, "git-receive-pack") {
// 			t := team{}
// 			err := first(map[string]interface{}{"name": teamname}, &t)
// 			if err != nil {
// 				respondError(w, http.StatusInternalServerError, err)
// 				return
// 			}
// 			username, password, ok := r.BasicAuth()
// 			u := user{}
// 			err = first(map[string]interface{}{"name": username}, &u)
// 			if err != nil {
// 				respondError(w, http.StatusInternalServerError, err)
// 				return
// 			}
// 			teamUser := teamUser{}
// 			err = first(map[string]interface{}{"team_id": t.ID, "user_id": u.ID}, &teamUser)
// 			if err != nil {
// 				respondError(w, http.StatusInternalServerError, err)
// 				return
// 			}
// 			if teamUser.ID == 0 {
// 				Unauthorized(w)
// 				return
// 			}
// 			err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
// 			if !ok || err != nil {
// 				Unauthorized(w)
// 				return
// 			}
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// func preflight(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodOptions {
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
