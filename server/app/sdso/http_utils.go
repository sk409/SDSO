package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

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
	db.Where("id = ?", userID).First(&u)
	if db.Error != nil {
		return nil, err
	}
	return &u, nil
}

func first(query map[string]interface{}, model interface{}) (int, error) {
	db.Where(query).First(model)
	if db.Error != nil {
		return http.StatusInternalServerError, db.Error
	}
	rv := reflect.ValueOf(model).Elem()
	id := rv.FieldByName("ID")
	if id.Uint() == 0 {
		return http.StatusBadRequest, errNotExist
	}
	return 0, nil
}

func fetch(r *http.Request, model interface{}) error {
	query := make(map[string]interface{})
	for key, value := range r.URL.Query() {
		query[key] = value[0]
	}
	db.Where(query).Find(model)
	return db.Error
}

func login(w http.ResponseWriter, name, password string) (*user, error) {
	u := user{}
	db.Where("name = ?", name).First(&u)
	if db.Error != nil {
		return nil, db.Error
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
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
	return &u, nil
}

func respond(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func respondError(w http.ResponseWriter, statusCode int, err error) {
	log.Println(err)
	respond(w, statusCode)
}

func respondJSON(w http.ResponseWriter, statusCode int, model interface{}) {
	jsonBytes, err := json.Marshal(model)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
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

func store(r *http.Request, model interface{}) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	rv := reflect.ValueOf(model).Elem()
	for key, value := range r.PostForm {
		fieldName := string(gocase.UpperCamelCase([]byte(key), true))
		fv := rv.FieldByName(fieldName)
		ft := fv.Type()
		if ft.Kind() == reflect.String {
			fv.SetString(value[0])
		} else if ft.Kind() == reflect.Uint {
			u, err := strconv.ParseUint(value[0], 10, 64)
			if err != nil {
				return http.StatusBadRequest, errBadRequest
			}
			fv.SetUint(u)
		}
	}
	db.Save(model)
	if db.Error != nil {
		return http.StatusInternalServerError, db.Error
	}
	return http.StatusOK, nil
}
