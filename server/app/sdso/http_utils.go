package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/sk409/gocase"

	"golang.org/x/crypto/bcrypt"

	"github.com/sk409/goconst"
)

func authenticatedUser(r *http.Request) (*user, error) {
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		return nil, err
	}
	sessionID := sessionCookie.Value
	session, err := sessionManager.Provider.Get(sessionID)
	if err != nil {
		return nil, err
	}
	userID, err := session.Uint(sessionStoreNameUserID)
	if err != nil {
		return nil, err
	}
	u, err := userRepository.findByID(userID, loadAllRelation)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func checkPermissionWithRequest(r *http.Request, users []user) (bool, error) {
	u, err := authenticatedUser(r)
	if err != nil {
		return false, err
	}
	return checkPermission(u, users), nil
}

func login(w http.ResponseWriter, username, password string) (*user, error) {
	u, err := userRepository.findByName(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	session, err := sessionManager.Provider.Start()
	if err != nil {
		return nil, err
	}
	session.Store(sessionStoreNameUserID, u.ID)
	sessionCookie := newCookie(cookieNameSessionID, session.ID(), cookie30Days)
	http.SetCookie(w, sessionCookie)
	return u, nil
}

func makeQuery(r *http.Request, model interface{}, snake bool) map[string]interface{} {
	rt := reflect.TypeOf(model)
	m := make(map[string]interface{})
	values := url.Values{}
	if r.Method == http.MethodGet {
		values = r.URL.Query()
	} else {
		r.ParseForm()
		values = r.PostForm
	}
	for key, value := range values {
		u := string(gocase.UpperCamelCase([]byte(key), true))
		if _, ok := rt.FieldByName(u); !ok {
			continue
		}
		k := key
		if snake {
			k = string(gocase.SnakeCase([]byte(key)))
		}
		m[k] = value[0]
	}
	return m
}

func respond(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func respondError(w http.ResponseWriter, statusCode int, err error) {
	log.Println(err)
	respond(w, statusCode)
}

func respondJSON(w http.ResponseWriter, statusCode int, model interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	// respond(w, statusCode)
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
	return jsonBytes, nil
}

func respondMessage(w http.ResponseWriter, statusCode int, message string) {
	respond(w, statusCode)
	w.Write([]byte(message))
}

func routeWithID(r *http.Request) (string, bool) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	components := strings.Split(path, "/")
	if len(components) < 4 {
		return "", false
	}
	return components[2], true
}

func putWebsocket(w http.ResponseWriter, r *http.Request, sockets *map[uint]*websocket.Conn) error {
	userID := r.URL.Query().Get("userId")
	if emptyAny(userID) {
		return errBadRequest
	}
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return errBadRequest
	}
	socket, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	(*sockets)[uint(userIDUint)] = socket
	return nil
}
