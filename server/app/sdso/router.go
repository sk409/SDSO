package main

import (
	"net/http"
	"path"
	"strings"

	"github.com/sk409/goconst"
	"golang.org/x/crypto/bcrypt"
)

type router struct {
	handler     http.Handler
	middlewares []func(w http.ResponseWriter, r *http.Request) bool
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var next bool
	for _, middleware := range router.middlewares {
		next = middleware(w, r)
		if !next {
			break
		}
	}
	if next {
		router.handler.ServeHTTP(w, r)
	}
}

func (router *router) allowCredentials(methods ...string) {
	router.middlewares = append(router.middlewares, func(w http.ResponseWriter, r *http.Request) bool {
		for _, method := range methods {
			if r.Method == method {
				w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
				break
			}
		}
		return true
	})
}

func (router *router) allowHeaders(methodHeaders map[string][]string) {
	router.middlewares = append(router.middlewares, func(w http.ResponseWriter, r *http.Request) bool {
		for method, headers := range methodHeaders {
			if r.Method == method {
				v := strings.Join(headers, ",")
				w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, v)
			}
		}
		return true
	})
}

func (router *router) allowMethods(methods ...string) {
	router.middlewares = append(router.middlewares, func(w http.ResponseWriter, r *http.Request) bool {
		v := strings.Join(methods, ",")
		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, v)
		return true
	})
}

func (router *router) cors() {
	router.middlewares = append(router.middlewares, func(w http.ResponseWriter, r *http.Request) bool {
		//TODO: Originをenvから取得
		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_ORIGIN, "http://localhost:3000")
		return true
	})
}

func (router *router) gitBasicAuth() {
	unauthorized := func(w http.ResponseWriter) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password."`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		w.Header().Set("Content-Type", "text/plain")
	}
	router.middlewares = append(router.middlewares, func(w http.ResponseWriter, r *http.Request) bool {
		components := strings.Split(r.URL.Path, "/")
		if len(components) < 3 {
			respond(w, http.StatusNotFound)
			return false
		}
		teamname := components[1]
		projectname := components[2]
		if r.URL.Path == "/"+path.Join(teamname, projectname, "git-receive-pack") {
			t := team{}
			err := first(map[string]interface{}{"name": teamname}, &t)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return false
			}
			username, password, ok := r.BasicAuth()
			u := user{}
			err = first(map[string]interface{}{"name": username}, &u)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return false
			}
			teamUser := teamUser{}
			err = first(map[string]interface{}{"team_id": t.ID, "user_id": u.ID}, &teamUser)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return false
			}
			if teamUser.ID == 0 {
				unauthorized(w)
				return false
			}
			err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
			if !ok || err != nil {
				unauthorized(w)
				return false
			}
		}
		return true
	})
}

func (router *router) preflight() {
	router.middlewares = append(router.middlewares, func(w http.ResponseWriter, r *http.Request) bool {
		if r.Method == http.MethodOptions {
			return false
		}
		return true
	})
}
