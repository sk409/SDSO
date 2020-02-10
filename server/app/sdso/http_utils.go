package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
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
	u := user{}
	gormDB.Where("id = ?", userID).First(&u)
	if gormDB.Error != nil {
		return nil, err
	}
	return &u, nil
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

func makeQuery(r *http.Request, snake bool) map[string]interface{} {
	m := make(map[string]interface{})
	values := url.Values{}
	if r.Method == http.MethodGet {
		values = r.URL.Query()
	} else {
		r.ParseForm()
		values = r.PostForm
	}
	for key, value := range values {
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

// func respondJSON(w http.ResponseWriter, statusCode int, model interface{}) ([]byte, error) {
// 	data := model
// 	data, err := public(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	jsonBytes, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
// 	w.Write(jsonBytes)
// 	return jsonBytes, nil
// 	//respond(w, statusCode)
// }

func respondJSON(w http.ResponseWriter, statusCode int, model interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
	return jsonBytes, nil
	//respond(w, statusCode)
}

func respondMessage(w http.ResponseWriter, statusCode int, message string) {
	w.Write([]byte(message))
	respond(w, statusCode)
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
