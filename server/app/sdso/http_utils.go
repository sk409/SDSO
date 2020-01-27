package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/sk409/gotype"

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

func fetch(r *http.Request, model interface{}) error {
	query := make(map[string]interface{})
	for key, value := range r.URL.Query() {
		s := string(gocase.SnakeCase([]byte(key)))
		query[s] = value[0]
	}
	db.Where(query).Find(model)
	return db.Error
}

func login(w http.ResponseWriter, username, password string) (*user, error) {
	u := user{}
	db.Where("name = ?", username).First(&u)
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
	data := model
	if gotype.IsSlice(data) {
		ft := reflect.TypeOf((*facade)(nil)).Elem()
		dt := reflect.TypeOf(data).Elem()
		if dt.Implements(ft) {
			dv := reflect.ValueOf(data)
			s := make([]interface{}, dv.Len())
			for i := 0; i < dv.Len(); i++ {
				p := dv.Index(i).Interface().(facade).public()
				s[i] = p
			}
			data = s
		}
	} else {
		if f, ok := model.(facade); ok {
			data = f.public()
		}
	}
	data, err := public(data)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	jsonBytes, err := json.Marshal(data)
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

func store(r *http.Request, model interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	query := make(map[string]interface{})
	for key, value := range r.PostForm {
		query[key] = value[0]
	}
	return save(query, model)
	// rv := reflect.ValueOf(model).Elem()
	// for key, value := range r.PostForm {
	// 	fieldName := string(gocase.UpperCamelCase([]byte(key), true))
	// 	fv := rv.FieldByName(fieldName)
	// 	ft := fv.Type()
	// 	if ft.Kind() == reflect.String {
	// 		fv.SetString(value[0])
	// 	} else if ft.Kind() == reflect.Uint {
	// 		u, err := strconv.ParseUint(value[0], 10, 64)
	// 		if err != nil {
	// 			return errBadRequest
	// 		}
	// 		fv.SetUint(u)
	// 	}
	// }
	// db.Save(model)
	// if db.Error != nil {
	// 	return db.Error
	// }
	// return nil
}
