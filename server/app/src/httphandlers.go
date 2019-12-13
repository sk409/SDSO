package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sk409/goconst"
)

func makeQueryAndValues(r *http.Request) (string, []interface{}) {
	conditions := []string{}
	values := []interface{}{}
	if r.Method == http.MethodGet {
		for key, value := range r.URL.Query() {
			conditions = append(conditions, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	} else {
		r.ParseForm()
		for key, value := range r.PostForm {
			conditions = append(conditions, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	}
	query := strings.Join(conditions, " AND ")
	return query, values
}

func existHandler(w http.ResponseWriter, r *http.Request, tableName string) {
	query, values := makeQueryAndValues(r)
	count := 0
	db.Table(tableName).Where(query, values...).Count(&count)
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonString := fmt.Sprintf("{\"exist\":%t}", count != 0)
	w.Write([]byte(jsonString))
}

func authCheckHandler(w http.ResponseWriter, r *http.Request) {
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodGet)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	notAuthenticated := []byte("{\"authenticated\": false}")
	authenticated := []byte("{\"authenticated\": true}")
	sessionCookie, err := r.Cookie("sessionID")
	if err != nil {
		w.Write(notAuthenticated)
		return
	}
	session, err := sessionManager.Provider.Get(sessionCookie.Value)
	if session == nil || err != nil {
		// log.Println(session)
		// if err != nil {
		// 	log.Println(err.Error())
		// }
		w.Write(notAuthenticated)
		http.SetCookie(w, &http.Cookie{
			Name:   cookieNameSessionID,
			MaxAge: -1,
		})
		return
	}
	w.Write(authenticated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodPost)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	if r.Method == http.MethodOptions {
		return
	}
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")
	if len(name) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(password)))
	user := user{}
	db.Model(user).Where("name = ? AND password = ?", name, hashedPassword).First(&user)
	exist := user.ID != 0
	if exist {
		session, err := sessionManager.Provider.Start()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		session.Store(sessionStoreNameUserID, user.ID)
		sessionCookie := newCookie(cookieNameSessionID, session.ID(), cookie30Days)
		http.SetCookie(w, sessionCookie)
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	response := fmt.Sprintf("{\"exist\": %t}", exist)
	w.Write([]byte(response))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodPost)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	if r.Method == http.MethodOptions {
		return
	}
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")
	if len(name) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(password)))
	user := user{Name: name, Password: hashedPassword}
	count := 0
	db.Model(user).Where("name = ? AND password = ?", user.Name, user.Password).Count(&count)
	if count != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	db.Save(&user)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session, err := sessionManager.Provider.Start()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.Store(sessionStoreNameUserID, user.ID)
	sessionCookie := newCookie(cookieNameSessionID, session.ID(), cookie30Days)
	http.SetCookie(w, sessionCookie)
}

func fetchUserHandler(w http.ResponseWriter, r *http.Request) {
	allowedHeaders := strings.Join([]string{goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_X_XSRF_TOKEN}, ",")
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, allowedHeaders)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodGet)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	if r.Method == http.MethodOptions {
		return
	}
	sessionCookie, err := r.Cookie(cookieNameSessionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session, err := sessionManager.Provider.Get(sessionCookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := session.Uint(sessionStoreNameUserID)
	user := user{}
	user.ID = userID
	db.Find(&user)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func fetchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query, values := makeQueryAndValues(r)
	users := []user{}
	db.Where(query, values...).Find(&users)
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func existUserHandler(w http.ResponseWriter, r *http.Request) {
	existHandler(w, r, tableNameUsers)
}

func storeProjectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.PostForm)
	name := r.PostFormValue("name")
	userIDString := r.PostFormValue("userID")
	if len(name) == 0 || len(userIDString) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db.Save(&project{Name: name, UserID: uint(userID)})
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func fetchProjectsHandler(w http.ResponseWriter, r *http.Request) {
	query, values := makeQueryAndValues(r)
	projects := []project{}
	db.Model(project{}).Where(query, values...).Find(&projects)
	jsonBytes, err := json.Marshal(projects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func existProjectHandler(w http.ResponseWriter, r *http.Request) {
	existHandler(w, r, tableNameProjects)
}

func storeScanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")
	projectName := r.PostFormValue("projectName")
	if len(userName) == 0 || len(password) == 0 || len(projectName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := user{}
	db.Where("name = ? AND password = ?", userName, password).First(&user)
	if db.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userIDString := strconv.Itoa(int(user.ID))
	project := project{}
	db.Where("name = ? AND user_id = ?", projectName, userIDString).First(&project)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	scan := scan{
		UserID:    user.ID,
		ProjectID: project.ID,
	}
	db.Save(&scan)
	jsonBytes, err := json.Marshal(scan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	w.Write(jsonBytes)
}

func fetchScansHandler(w http.ResponseWriter, r *http.Request) {
	query, values := makeQueryAndValues(r)
	scans := []scan{}
	db.Where(query, values...).Find(&scans)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonBytes, err := json.Marshal(scans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func storeVulnerability(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_METHODS, http.MethodPost)
	w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, goconst.HTTP_HEADER_CONTENT_TYPE)
	if r.Method == http.MethodOptions {
		return
	}
	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	path := r.PostFormValue("path")
	method := r.PostFormValue("method")
	request := r.PostFormValue("request")
	response := r.PostFormValue("response")
	scanID := r.PostFormValue("scanID")
	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")
	projectName := r.PostFormValue("projectName")
	if len(name) == 0 ||
		len(description) == 0 ||
		len(path) == 0 ||
		len(method) == 0 ||
		len(request) == 0 ||
		len(response) == 0 ||
		len(scanID) == 0 ||
		len(userName) == 0 ||
		len(password) == 0 ||
		len(projectName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	scanIDInt, err := strconv.Atoi(scanID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := &user{}
	db.Model(user).Where("name = ? AND password = ?", userName, password).First(&user)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	project := &project{}
	db.Model(project).Where("name = ? AND user_id = ?", projectName, user.ID).First(&project)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vulnerability := vulnerability{
		Name:        name,
		Description: description,
		Path:        path,
		Method:      method,
		Request:     request,
		Response:    response,
		ProjectID:   project.ID,
		ScanID:      uint(scanIDInt),
	}
	db.Save(&vulnerability)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func fetchVulnerabilities(w http.ResponseWriter, r *http.Request) {
	query, values := makeQueryAndValues(r)
	vulnerabilities := []vulnerability{}
	db.Where(query, values...).Find(&vulnerabilities)
	if db.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_JSON)
	jsonBytes, err := json.Marshal(vulnerabilities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}
